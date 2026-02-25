package sshserver

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"ssh-microvm/internal/config"
	"ssh-microvm/internal/firecracker"
	"ssh-microvm/internal/store"
	"ssh-microvm/internal/util"
)

type Server struct {
	cfg        *config.Config
	store      *store.Store
	manager    *firecracker.Manager
	hostSigner ssh.Signer
}

func New(cfg *config.Config, st *store.Store, manager *firecracker.Manager) (*Server, error) {
	signer, err := loadOrCreateHostKey(cfg.HostKeyPath)
	if err != nil {
		return nil, err
	}
	return &Server{cfg: cfg, store: st, manager: manager, hostSigner: signer}, nil
}

func (s *Server) Serve(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.cfg.ListenAddr)
	if err != nil {
		return err
	}
	return s.ServeListener(ctx, ln)
}

func (s *Server) ServeListener(ctx context.Context, ln net.Listener) error {
	defer ln.Close()

	go func() {
		<-ctx.Done()
		_ = ln.Close()
	}()

	log.Printf("ssh listening on %s", ln.Addr().String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			log.Printf("accept error: %v", err)
			continue
		}
		go s.handleConn(ctx, conn)
	}
}

func (s *Server) handleConn(ctx context.Context, netConn net.Conn) {
	serverCfg := &ssh.ServerConfig{
		PublicKeyCallback: s.publicKeyCallback,
	}
	serverCfg.AddHostKey(s.hostSigner)

	sshConn, channels, requests, err := ssh.NewServerConn(netConn, serverCfg)
	if err != nil {
		_ = netConn.Close()
		return
	}
	defer sshConn.Close()

	go ssh.DiscardRequests(requests)

	perm := sshConn.Permissions
	fingerprint := ""
	publicKey := ""
	if perm != nil && perm.Extensions != nil {
		fingerprint = perm.Extensions["pubkey-fp"]
		publicKey = perm.Extensions["pubkey"]
	}

	if fingerprint == "" || publicKey == "" {
		_ = sshConn.Close()
		return
	}

	userID, err := s.store.EnsureUserAndKey(ctx, sshConn.User(), fingerprint, publicKey)
	if err != nil {
		log.Printf("db ensure user/key error: %v", err)
		_ = sshConn.Close()
		return
	}

	sessionID, err := util.RandomHex(8)
	if err != nil {
		_ = sshConn.Close()
		return
	}

	session := store.Session{
		ID:             sessionID,
		UserID:         userID,
		KeyFingerprint: fingerprint,
		RemoteAddr:     sshConn.RemoteAddr().String(),
		StartedAt:      time.Now().UTC().Format(time.RFC3339Nano),
		Status:         "active",
	}
	if err := s.store.CreateSession(ctx, session); err != nil {
		log.Printf("db create session error: %v", err)
		_ = sshConn.Close()
		return
	}

	vm, err := s.manager.Start(ctx)
	if err != nil {
		_ = s.store.EndSession(ctx, sessionID, "vm_failed")
		log.Printf("firecracker start error: %v", err)
		_ = sshConn.Close()
		return
	}

	vmRecord := store.VM{
		ID:        vm.ID,
		SessionID: sessionID,
		StateDir:  vm.StateDir,
		FCPid:     vm.Cmd.Process.Pid,
		StartedAt: time.Now().UTC().Format(time.RFC3339Nano),
	}
	if err := s.store.CreateVM(ctx, vmRecord); err != nil {
		log.Printf("db create vm error: %v", err)
	}
	if err := s.store.AttachVM(ctx, sessionID, vm.ID); err != nil {
		log.Printf("db attach vm error: %v", err)
	}

	go s.handleChannels(channels, vm)

	err = sshConn.Wait()
	if err != nil {
		log.Printf("ssh connection ended: %v", err)
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.GracefulStopS)*time.Second)
	defer cancel()
	_ = vm.Stop(stopCtx, time.Duration(s.cfg.GracefulStopS)*time.Second)
	_ = s.store.EndVM(ctx, vm.ID, 0)
	_ = s.store.EndSession(ctx, sessionID, "closed")
}

func (s *Server) handleChannels(channels <-chan ssh.NewChannel, vm *firecracker.VM) {
	for ch := range channels {
		if ch.ChannelType() != "session" {
			_ = ch.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := ch.Accept()
		if err != nil {
			continue
		}
		go s.handleSession(channel, requests, vm)
	}
}

func (s *Server) handleSession(ch ssh.Channel, requests <-chan *ssh.Request, vm *firecracker.VM) {
	defer ch.Close()

	var (
		ptyReq    *ptyRequest
		startErr  error
		startOnce sync.Once
		winCh     = make(chan windowChange, 8)
	)

	startGuest := func(shell bool, execCmd string) {
		startOnce.Do(func() {
			startErr = s.proxyToGuest(ch, ptyReq, winCh, shell, execCmd, vm)
			if startErr != nil {
				log.Printf("guest proxy error: %v", startErr)
			}
		})
	}

	for req := range requests {
		switch req.Type {
		case "pty-req":
			pty, ok := parsePtyRequest(req.Payload)
			if ok {
				ptyReq = &pty
				_ = req.Reply(true, nil)
			} else {
				_ = req.Reply(false, nil)
			}
		case "window-change":
			if win, ok := parseWindowChange(req.Payload); ok {
				select {
				case winCh <- win:
				default:
				}
			}
		case "shell":
			_ = req.Reply(true, nil)
			startGuest(true, "")
		case "exec":
			cmd, ok := parseExecRequest(req.Payload)
			if ok {
				_ = req.Reply(true, nil)
				startGuest(false, cmd)
			} else {
				_ = req.Reply(false, nil)
			}
		case "env":
			_ = req.Reply(true, nil)
		default:
			_ = req.Reply(false, nil)
		}
	}
}

func (s *Server) publicKeyCallback(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	fingerprint := ssh.FingerprintSHA256(key)
	pubKey := string(ssh.MarshalAuthorizedKey(key))

	if s.cfg.AuthMode == config.AuthModeKnownKeys {
		exists, err := s.store.HasKey(context.Background(), fingerprint)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("unknown key: %s", fingerprint)
		}
	}

	return &ssh.Permissions{
		Extensions: map[string]string{
			"pubkey-fp": fingerprint,
			"pubkey":    pubKey,
		},
	}, nil
}

func loadOrCreateHostKey(path string) (ssh.Signer, error) {
	if data, err := os.ReadFile(path); err == nil {
		return ssh.ParsePrivateKey(data)
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}

	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	if err := os.WriteFile(path, pemBytes, 0o600); err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(pemBytes)
}

type ptyRequest struct {
	Term   string
	Width  int
	Height int
}

type windowChange struct {
	Width  int
	Height int
}

func parsePtyRequest(payload []byte) (ptyRequest, bool) {
	var req struct {
		Term          string
		Columns, Rows uint32
		Width, Height uint32
		TerminalModes []byte
	}
	if err := ssh.Unmarshal(payload, &req); err != nil {
		return ptyRequest{}, false
	}
	return ptyRequest{
		Term:   req.Term,
		Width:  int(req.Columns),
		Height: int(req.Rows),
	}, true
}

func parseWindowChange(payload []byte) (windowChange, bool) {
	var req struct {
		Columns, Rows uint32
		Width, Height uint32
	}
	if err := ssh.Unmarshal(payload, &req); err != nil {
		return windowChange{}, false
	}
	return windowChange{
		Width:  int(req.Columns),
		Height: int(req.Rows),
	}, true
}

func parseExecRequest(payload []byte) (string, bool) {
	var req struct {
		Command string
	}
	if err := ssh.Unmarshal(payload, &req); err != nil {
		return "", false
	}
	return req.Command, true
}

func (s *Server) proxyToGuest(ch ssh.Channel, ptyReq *ptyRequest, winCh <-chan windowChange, shell bool, execCmd string, vm *firecracker.VM) error {
	if vm == nil {
		return errors.New("vm not available")
	}
	if err := waitForPort(fmt.Sprintf("%s:22", vm.GuestIP), 15*time.Second); err != nil {
		return err
	}

	client, err := s.dialGuest(vm.GuestIP)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if ptyReq != nil {
		if err := session.RequestPty(ptyReq.Term, ptyReq.Height, ptyReq.Width, ssh.TerminalModes{}); err != nil {
			return err
		}
	}

	session.Stdin = ch
	session.Stdout = ch
	session.Stderr = ch

	go func() {
		for win := range winCh {
			_ = session.WindowChange(win.Height, win.Width)
		}
	}()

	if shell {
		if err := session.Shell(); err != nil {
			return err
		}
		return session.Wait()
	}

	if execCmd == "" {
		execCmd = "bash"
	}
	return session.Run(execCmd)
}

func (s *Server) dialGuest(guestIP string) (*ssh.Client, error) {
	keyData, err := os.ReadFile(s.cfg.GuestKeyPath)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, err
	}
	clientCfg := &ssh.ClientConfig{
		User:            s.cfg.GuestUser,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	return ssh.Dial("tcp", fmt.Sprintf("%s:22", guestIP), clientCfg)
}

func waitForPort(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for %s", addr)
}
