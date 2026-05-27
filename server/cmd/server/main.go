package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tyzerrr/aws-log-practice/server/gen/greet/v1/v1connect"
	"github.com/tyzerrr/aws-log-practice/server/internal/adapter/db"
)

const (
	ExitOK int = iota
	ExitErr
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		os.Exit(ExitErr)
	}
	os.Exit(ExitOK)
}

func run(logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// db pool
	dbPool, err := db.NewDBPool(ctx, logger, os.Getenv("DB_URL"))
	if err != nil {
		logger.Error("failed to start db pool", slog.String("error", err.Error()))
		return err
	}
	defer dbPool.Close()

	// http server
	addr := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	mux := http.NewServeMux()
	path, greetHandler := v1connect.NewGreetServiceHandler(NewGreetHandler())
	mux.Handle(path, greetHandler)

	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)

	srv := http.Server{
		Addr:      addr,
		Handler:   mux,
		Protocols: p,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	logger.Info("server starting", slog.String("addr", addr))
	select {
	case <-ctx.Done():
		logger.Info("server start to shutdown...")
	case err := <-errCh:
		if err != nil {
			logger.Error("server failed", slog.String("error", err.Error()))
			return err
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server failed to shutdown gracefully", slog.String("error", err.Error()))
		return err
	}
	logger.Info("server success to shutdown!")
	return nil
}
