package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/thetnaingtn/dirty-hand/internal/config"
	"github.com/thetnaingtn/dirty-hand/server"
	"github.com/thetnaingtn/dirty-hand/store"
	"github.com/thetnaingtn/dirty-hand/store/db"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load configuration using the new Viper-based system with --config flag support
	config, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	driver, err := db.NewDBDriver(config)
	if err != nil {
		slog.Error("failed to create database driver", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	store := store.NewStore(driver, config)

	s, err := server.NewServer(ctx, store, config)
	if err != nil {
		cancel()
		slog.Error("failed to create server", "error", err)
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	if err := s.Start(); err != nil {
		if err != http.ErrServerClosed {
			cancel()
			slog.Error("failed to start server", "error", err)
		}
	}

	slog.Info("Server started", "address", fmt.Sprintf("%s:%s", config.Server.Addr, config.Server.Port), "environment", config.Environment)

	go func() {
		<-c
		if err := s.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown server", "error", err)
		}
		cancel()
	}()

	<-ctx.Done()
}
