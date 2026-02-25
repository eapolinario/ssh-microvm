package main

import (
	"context"
	"log"
	"os"

	"ssh-microvm/internal/config"
	"ssh-microvm/internal/firecracker"
	"ssh-microvm/internal/sshserver"
	"ssh-microvm/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	if err := os.MkdirAll(cfg.StateDir, 0o750); err != nil {
		log.Fatalf("state dir error: %v", err)
	}

	st, err := store.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("db open error: %v", err)
	}
	defer st.Close()

	if err := st.EnsureSchema(context.Background()); err != nil {
		log.Fatalf("db schema error: %v", err)
	}

	manager := firecracker.NewManager(cfg)
	server, err := sshserver.New(cfg, st, manager)
	if err != nil {
		log.Fatalf("ssh server init error: %v", err)
	}

	if err := server.Serve(context.Background()); err != nil {
		log.Fatalf("ssh server error: %v", err)
	}
}
