package main

import (
	"log/slog"
	"longtasks/internal/app"
	"longtasks/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	log.Info("Starting...")

	application := app.New(log, cfg)

	go application.HTTPServer.MustRun()

	//graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	log.Info("Stopping application")

	application.Storage.Close()
	application.HTTPServer.Stop()

	log.Info("Application stopped")
}
