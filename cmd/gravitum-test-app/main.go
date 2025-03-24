package main

import (
	"context"
	"fmt"
	"gravitum-test-app/config"
	"gravitum-test-app/internal/app"
	"gravitum-test-app/pkg/logger"
	slog "log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	defer fmt.Println("Server shutdown")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		slog.Fatalf("config error: %s", err)
	}

	cfg.Print()

	log := logger.New(logger.GetLevelByString(cfg.Log.Level))
	a := app.New(&cfg, log)

	// Channel to listen for OS signals
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Run the app in a goroutine
	go func() {
		if err := a.Run(ctx); err != nil {
			log.Error(fmt.Sprintf("app run: %s", err))
		}
	}()

	// Wait for shutdown signal
	<-osSignal
	log.Info("Received shutdown signal, waiting for ongoing transactions to complete...")

	// signal all services using ctx to stop
	cancel()

	time.Sleep(1 * time.Second)

	log.Info("terminating server")
	a.Server.Shutdown(context.Background())
}
