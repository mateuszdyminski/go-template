package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mateuszdyminski/go-template/repository/postgres"

	"go.uber.org/zap"
)

// @title Go-template
// @description This is a simple http service written in go.

// @contact.name Mateusz Dyminski
// @contact.email dyminski@gmail.com

// @BasePath /
func main() {
	logger := initLogger()
	ls := logger.Sugar()

	cfg, err := loadConfig(logger)
	if err != nil {
		ls.Fatalw("can't load configuration", "err", err)
	}

	repo, err := postgres.NewPostgresRepository(cfg.pgHost, cfg.pgPort, cfg.pgUser, cfg.pgPassword, cfg.pgDBName)
	if err != nil {
		ls.Fatalw("can't create repository", "err", err)
	}

	// wait for SIGTERM or SIGINT
	cancelCtx := initContext()

	router := newRouter(cancelCtx, logger, repo)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.httpPort),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 1 * time.Minute,
		IdleTimeout:  15 * time.Second,
	}

	// run server in background
	go func() {
		ls.Infow("HTTP Server started", "port", cfg.httpPort)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			ls.Fatalw("can't start HTTP server", "err", err)
		}
	}()

	// run pprof server in background on different port
	if cfg.httpPprofPort != 0 {
		go func() {
			ls.Infow("HTTP Server for pprof purpose started", "port", cfg.httpPprofPort)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.httpPprofPort), newPprofRouter()); err != nil {
				ls.Fatalw("can't start pprof HTTP server", "err", err)
			}
		}()
	}

	<-cancelCtx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.httpGracefulTimeout)*time.Second)
	defer cancel()

	// wait for Kubernetes readiness probe to remove this instance from service
	// the readiness check interval must be lower than the timeout

	// sleep some additional time to drain all ongoing requests - use it only on prod
	// as long as it's annoing on development
	if !debug() {
		time.Sleep(time.Duration(int64(cfg.httpGracefulSleep) * int64(time.Second)))
	}

	ls.Infow("shutting down HTTP server", "timeout", time.Duration(cfg.httpGracefulTimeout)*time.Second)
	srv.SetKeepAlivesEnabled(false)

	if err := srv.Shutdown(ctx); err != nil {
		ls.Errorw("HTTP server graceful shutdown failed", "err", err)
	} else {
		ls.Infow("HTTP server gracefully stopped")
	}
}

func initContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}

func initLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	if debug() {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("can't init logger: %s", err)
	}

	return logger
}

func debug() bool {
	d := os.Getenv("DEBUG")
	if d == "1" || d == "true" || d == "True" {
		return true
	}

	return false
}
