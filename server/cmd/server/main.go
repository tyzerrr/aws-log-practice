package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/tyzerrr/aws-log-practice/server/cmd/server/handler"
	productv1connect "github.com/tyzerrr/aws-log-practice/server/gen/product/v1/v1connect"
	"github.com/tyzerrr/aws-log-practice/server/internal/adapter/db"
	"github.com/tyzerrr/aws-log-practice/server/internal/infrastructure"
	"github.com/tyzerrr/aws-log-practice/server/internal/usecase"
)

const (
	ExitOK int = iota
	ExitErr
)

const (
	defaultDBPort    = "5432"
	defaultDBSSLMode = "require"
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
	databaseURL, err := databaseURLFromEnv()
	if err != nil {
		logger.Error("failed to load database config", slog.String("error", err.Error()))
		return err
	}
	dbPool, err := db.NewDBPool(ctx, logger, databaseURL)
	if err != nil {
		logger.Error("failed to start db pool", slog.String("error", err.Error()))
		return err
	}
	defer dbPool.Close()

	txManager := infrastructure.NewTransactionManager(dbPool.Pool)

	// http server
	addr := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	mux := http.NewServeMux()

	// products
	productUsecase := usecase.NewProductUsecase(
		logger,
		infrastructure.NewProductRepository,
		txManager,
	)
	productPath, productHandler := productv1connect.NewProductServiceHandler(
		handler.NewProductHandler(
			logger,
			productUsecase,
		),
	)
	mux.Handle(productPath, productHandler)

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

func databaseURLFromEnv() (string, error) {
	if databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL")); databaseURL != "" {
		return databaseURL, nil
	}

	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	port := envDefault("DB_PORT", defaultDBPort)
	name := strings.TrimSpace(os.Getenv("DB_NAME"))
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	password := strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	sslMode := envDefault("DB_SSLMODE", defaultDBSSLMode)

	switch {
	case host == "":
		return "", fmt.Errorf("DB_HOST is required when DATABASE_URL is empty")
	case name == "":
		return "", fmt.Errorf("DB_NAME is required when DATABASE_URL is empty")
	case user == "":
		return "", fmt.Errorf("DB_USER is required when DATABASE_URL is empty")
	case password == "":
		return "", fmt.Errorf("DB_PASSWORD is required when DATABASE_URL is empty")
	case sslMode == "":
		return "", fmt.Errorf("DB_SSLMODE is required when DATABASE_URL is empty")
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   "/" + name,
	}
	values := url.Values{}
	values.Set("sslmode", sslMode)
	u.RawQuery = values.Encode()

	return u.String(), nil
}

func envDefault(key string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}
