package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testTask/internal/config"
	"testTask/internal/question"
	"testTask/internal/question/db"
	"testTask/pkg/client/postgres"
	"testTask/pkg/logging"
	"time"
)

func main() {
	logger := logging.GetLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("config error: %v", err)
	}

	ctx := context.Background()

	client, err := postgres.NewClient(ctx, cfg.DSN)
	if err != nil {
		logger.Fatalf("postgres init error: %v", err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			logger.Warnf("postgres close error: %v", err)
		}
	}()

	mux := http.NewServeMux()

	questionStorage := db.NewStorage(client.DB, logger)
	questionService := question.NewService(questionStorage, logger)
	questionHandler := question.NewHandler(logger, questionService)
	questionHandler.Register(mux)

	answerStorage := db.NewStorage(client.DB, logger)
	answerService := question.NewService(answerStorage, logger)
	answerHandler := question.NewHandler(logger, answerService)
	answerHandler.Register(mux)

	startServer(mux)
}

func startServer(mux *http.ServeMux) {
	logger := logging.GetLogger()

	srv := &http.Server{Addr: os.Getenv("APP_PORT"), Handler: mux}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen: %v", err)
		}
	}()

	logger.Info("server started")
	<-ctx.Done()
	logger.Info("server stopping...")

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shCtx); err != nil {
		logger.Fatalf("shutdown failed: %v", err)
	}

	logger.Info("server stopped")
}
