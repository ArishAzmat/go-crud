package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arishazmat/go-crud/internal/config"
	todo "github.com/arishazmat/go-crud/internal/http/handlers/students"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// connect to db
	// set router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/todos", todo.New())

	// server setup

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Server is starting...", slog.String("address", cfg.Addr))
	// fmt.Printf("Server is running on url: %s", cfg.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server: ", err)
		}
	}()

	<-done

	slog.Info("server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	/* if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server: ", slog.String("error", err.Error()))
	}
	*/
	if err != nil {
		slog.Error("Failed to shutdown server: ", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown gracefully")

}
