package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store, err := NewMySQLStore(cfg.MySQLDSN, cfg.MySQLMaxIdle, cfg.MySQLMaxOpen, cfg.MySQLMaxLifeMin)
	if err != nil {
		panic(err)
	}
	defer func() { _ = store.Close() }()

	if err = store.EnsureSchema(ctx); err != nil {
		panic(err)
	}

	client := NewActivityClient(cfg.DataAPIBaseURL, cfg.TargetUser, cfg.HTTPTimeout)
	poller := NewPoller(cfg, client, store)
	go poller.Run(ctx)

	h := NewHTTPServer(cfg, poller, store)
	errCh := make(chan error, 1)
	go func() {
		log.Printf("[pm_activity] server start addr=%s user=%s limit=%d fast=%s backfill=%s offsets=%v", cfg.ListenAddr, cfg.TargetUser, cfg.PageLimit, cfg.FastInterval, cfg.BackfillInterval, cfg.BackfillOffsets)
		errCh <- h.Run()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case sig := <-sigCh:
		log.Printf("[pm_activity] received signal=%v, shutting down", sig)
	case err = <-errCh:
		if err != nil {
			log.Printf("[pm_activity] server exited with error: %v", err)
		}
		cancel()
		return
	}

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer shutdownCancel()
	if err = h.Shutdown(shutdownCtx); err != nil {
		log.Printf("[pm_activity] graceful shutdown failed: %v", err)
	}
	time.Sleep(400 * time.Millisecond)
	if err != nil {
		panic(fmt.Errorf("server exit with err: %w", err))
	}
}
