package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testTask/internal/question"
	"testTask/pkg/logging"
	"time"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	name := r.PathValue("name")
	fmt.Fprintf(w, "Hello, %s", name)
}

func main() {
	//mux.HandleFunc("GET /{name}", IndexHandler)
	logger := logging.GetLogger()

	mux := http.NewServeMux()

	handler := question.NewHandler(logger)
	handler.Register(mux)

	startServer(mux)
}

func startServer(mux *http.ServeMux) {
	logger := logging.GetLogger()

	srv := &http.Server{Addr: ":8080", Handler: mux}
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
