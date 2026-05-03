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
	"reflect"
	"strings"
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
	if cfg == nil {
		return nil, errors.New("config must be set")
	}
	if st == nil {
		return nil, errors.New("store must be set")
	}
	if manager == nil {
		return nil, errors.New("firecracker manager must be set")
	}
	signer, err := loadOrCreateHostKey(cfg.HostKeyPath)
	if err != nil {
		return nil, err
	}
	return &Server{cfg: cfg, store: st, manager: manager, hostSigner: signer}, nil
}

func (s *Server) Serve(ctx context.Context) error {
	if err := s.validateServe(ctx); err != nil {
		return err
	}
	ln, err := net.Listen("tcp", s.cfg.ListenAddr)
	if err != nil {
		return err
	}
	return s.ServeListener(ctx, ln)
}

func (s *Server) ServeListener(ctx context.Context, ln net.Listener) error {
	if err := s.validateServe(ctx); err != nil {
		return err
	}
	if ln == nil {
		return errors.New("listener must be set")
	}
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

func (s *Server) validateServe(ctx context.Context) error {
	if s == nil {
		return errors.New("server must be set")
	}
	if ctx == nil {
		return errors.New("context must be set")
	}
	if s.cfg == nil {
		return errors.New("config must be set")
	}
	if s.store == nil {
		return errors.New("store must be set")
	}
	if s.manager == nil {
		return errors.New("firecracker manager must be set")
	}
	if s.hostSigner == nil {
		return errors.New("host signer must be set")
	}
	return nil
}

func (s *Server) handleConn(ctx context.Context, netConn net.Conn) {
	if err := s.validateServe(ctx); err != nil {
		if netConn != nil {
			_ = netConn.Close()
		}
		log.Printf("ssh connection rejected: %v", err)
		return
	}
	if netConn == nil {
		log.Printf("ssh connection rejected: network connection must be set")
		return
	}

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

	if err := s.stopVM(vm); err != nil {
		log.Printf("firecracker stop error: %v", err)
	}
	_ = s.store.EndVM(ctx, vm.ID, 0)
	_ = s.store.EndSession(ctx, sessionID, "closed")
}

func (s *Server) stopVM(vm *firecracker.VM) error {
	if s == nil {
		return errors.New("server must be set")
	}
	if s.cfg == nil {
		return errors.New("config must be set")
	}
	if vm == nil {
		return errors.New("vm not available")
	}
	if s.cfg.GracefulStopS <= 0 {
		return errors.New("graceful shutdown timeout must be > 0")
	}
	graceful := time.Duration(s.cfg.GracefulStopS) * time.Second
	stopCtx, cancel := context.WithTimeout(context.Background(), graceful)
	defer cancel()
	return vm.Stop(stopCtx, graceful)
}

func (s *Server) handleChannels(channels <-chan ssh.NewChannel, vm *firecracker.VM) {
	if s == nil {
		log.Printf("ssh channels rejected: server must be set")
		return
	}
	if channels == nil {
		log.Printf("ssh channels rejected: channel stream must be set")
		return
	}
	if vm == nil {
		log.Printf("ssh channels rejected: vm not available")
		return
	}

	for ch := range channels {
		if isNilSSHNewChannel(ch) {
			log.Printf("ssh channel rejected: new channel must be set")
			continue
		}
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
	if isNilSSHChannel(ch) {
		log.Printf("ssh session rejected: channel must be set")
		return
	}
	defer ch.Close()
	if s == nil {
		log.Printf("ssh session rejected: server must be set")
		return
	}
	if s.cfg == nil {
		log.Printf("ssh session rejected: config must be set")
		return
	}
	if requests == nil {
		log.Printf("ssh session rejected: requests must be set")
		return
	}
	if vm == nil {
		log.Printf("ssh session rejected: vm not available")
		return
	}

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
		if req == nil {
			log.Printf("ssh session request rejected: request must be set")
			continue
		}
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

func isNilSSHChannel(ch ssh.Channel) bool {
	return isNilInterface(ch)
}

func isNilSSHNewChannel(ch ssh.NewChannel) bool {
	return isNilInterface(ch)
}

func isNilInterface(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

func (s *Server) publicKeyCallback(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	if s == nil {
		return nil, errors.New("server must be set")
	}
	if s.cfg == nil {
		return nil, errors.New("config must be set")
	}
	if key == nil {
		return nil, errors.New("public key must be set")
	}

	fingerprint := ssh.FingerprintSHA256(key)
	pubKey := string(ssh.MarshalAuthorizedKey(key))

	switch s.cfg.AuthMode {
	case config.AuthModeAutoEnroll:
	case config.AuthModeKnownKeys:
		if s.store == nil {
			return nil, errors.New("store must be set")
		}
		exists, err := s.store.HasKey(context.Background(), fingerprint)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("unknown key: %s", fingerprint)
		}
	default:
		return nil, fmt.Errorf("invalid auth mode: %s", s.cfg.AuthMode)
	}

	return &ssh.Permissions{
		Extensions: map[string]string{
			"pubkey-fp": fingerprint,
			"pubkey":    pubKey,
		},
	}, nil
}

func loadOrCreateHostKey(path string) (ssh.Signer, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("host key path must be set")
	}
	if path != strings.TrimSpace(path) {
		return nil, errors.New("host key path must not contain surrounding whitespace")
	}
	if err := ensureHostKeyDir(filepath.Dir(path)); err != nil {
		return nil, err
	}
	if data, err := os.ReadFile(path); err == nil {
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if info.Mode().Perm()&0o077 != 0 {
			return nil, fmt.Errorf("host key %s permissions too open: %v", path, info.Mode().Perm())
		}
		return ssh.ParsePrivateKey(data)
	} else if !errors.Is(err, os.ErrNotExist) {
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

func ensureHostKeyDir(path string) error {
	if err := os.MkdirAll(path, 0o700); err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("host key directory %s is not a directory", path)
	}
	if info.Mode().Perm()&0o022 != 0 {
		return fmt.Errorf("host key directory %s permissions too open: %v", path, info.Mode().Perm())
	}
	return nil
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
	if strings.TrimSpace(req.Term) == "" || req.Columns == 0 || req.Rows == 0 {
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
	if req.Columns == 0 || req.Rows == 0 {
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
	if strings.TrimSpace(req.Command) == "" {
		return "", false
	}
	return req.Command, true
}

func (s *Server) proxyToGuest(ch ssh.Channel, ptyReq *ptyRequest, winCh <-chan windowChange, shell bool, execCmd string, vm *firecracker.VM) error {
	if err := s.validateGuestProxy(ch, vm); err != nil {
		return err
	}
	if err := validateGuestPTY(ptyReq); err != nil {
		return err
	}
	if err := validateWindowChanges(winCh); err != nil {
		return err
	}
	if err := validateGuestCommand(shell, execCmd); err != nil {
		return err
	}
	if err := s.validateGuestDial(vm.GuestIP); err != nil {
		return err
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
	return session.Run(execCmd)
}

func (s *Server) validateGuestProxy(ch ssh.Channel, vm *firecracker.VM) error {
	if s == nil {
		return errors.New("server must be set")
	}
	if s.cfg == nil {
		return errors.New("config must be set")
	}
	if isNilSSHChannel(ch) {
		return errors.New("ssh channel must be set")
	}
	if vm == nil {
		return errors.New("vm not available")
	}
	return nil
}

func validateGuestPTY(ptyReq *ptyRequest) error {
	if ptyReq == nil {
		return nil
	}
	if strings.TrimSpace(ptyReq.Term) == "" {
		return errors.New("pty terminal must be set")
	}
	if ptyReq.Width <= 0 || ptyReq.Height <= 0 {
		return errors.New("pty dimensions must be positive")
	}
	return nil
}

func validateWindowChanges(winCh <-chan windowChange) error {
	if winCh == nil {
		return errors.New("window change channel must be set")
	}
	return nil
}

func validateGuestCommand(shell bool, execCmd string) error {
	if shell {
		if strings.TrimSpace(execCmd) != "" {
			return errors.New("exec command cannot be set for shell sessions")
		}
		return nil
	}
	if strings.TrimSpace(execCmd) == "" {
		return errors.New("exec command must be set")
	}
	return nil
}

func (s *Server) validateGuestDial(guestIP string) error {
	if s == nil {
		return errors.New("server must be set")
	}
	if s.cfg == nil {
		return errors.New("config must be set")
	}
	if strings.TrimSpace(guestIP) == "" {
		return errors.New("guest IP must be set")
	}
	if !isIPv4(guestIP) {
		return fmt.Errorf("guest IP must be a valid IPv4 address: %s", guestIP)
	}
	if strings.TrimSpace(s.cfg.GuestUser) == "" {
		return errors.New("guest user must be set")
	}
	if s.cfg.GuestUser != strings.TrimSpace(s.cfg.GuestUser) {
		return errors.New("guest user must not contain surrounding whitespace")
	}
	if strings.TrimSpace(s.cfg.GuestKeyPath) == "" {
		return errors.New("guest key path must be set")
	}
	if s.cfg.GuestKeyPath != strings.TrimSpace(s.cfg.GuestKeyPath) {
		return errors.New("guest key path must not contain surrounding whitespace")
	}
	return nil
}

func (s *Server) dialGuest(guestIP string) (*ssh.Client, error) {
	if err := s.validateGuestDial(guestIP); err != nil {
		return nil, err
	}
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
	return waitForPortWithDial(addr, timeout, func(addr string, timeout time.Duration) (net.Conn, error) {
		return net.DialTimeout("tcp", addr, timeout)
	})
}

func isIPv4(value string) bool {
	ip := net.ParseIP(value)
	ipv4 := ip.To4()
	return ipv4 != nil && value == ipv4.String()
}

func waitForPortWithDial(addr string, timeout time.Duration, dial func(string, time.Duration) (net.Conn, error)) error {
	if strings.TrimSpace(addr) == "" {
		return errors.New("guest port address must be set")
	}
	if addr != strings.TrimSpace(addr) {
		return errors.New("guest port address must not contain surrounding whitespace")
	}
	if _, port, err := net.SplitHostPort(addr); err != nil {
		return fmt.Errorf("guest port address must be a valid TCP address: %s", addr)
	} else if port == "" {
		return errors.New("guest port address port must be set")
	} else if _, err := net.LookupPort("tcp", port); err != nil {
		return fmt.Errorf("guest port address port must be valid: %s", port)
	}
	if timeout <= 0 {
		return errors.New("guest port timeout must be positive")
	}
	if dial == nil {
		return errors.New("guest port dial function must be set")
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		dialTimeout := time.Until(deadline)
		if dialTimeout > 500*time.Millisecond {
			dialTimeout = 500 * time.Millisecond
		}
		conn, err := dial(addr, dialTimeout)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		sleepFor := time.Until(deadline)
		if sleepFor > 200*time.Millisecond {
			sleepFor = 200 * time.Millisecond
		}
		if sleepFor > 0 {
			time.Sleep(sleepFor)
		}
	}
	return fmt.Errorf("timeout waiting for %s", addr)
}
