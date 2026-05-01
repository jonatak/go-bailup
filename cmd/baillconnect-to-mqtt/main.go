package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonatak/baillconnect-to-mqtt/internal/bootstrap"
	"github.com/jonatak/baillconnect-to-mqtt/internal/config"
)

func main() {
	slog.Info(fmt.Sprintf("Start baillconnect-to-mqtt version:%s, commit:%s, buildtime:%s", config.Version, config.CommitSHA, config.BuildTime))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	service, err := bootstrap.NewHVACService()

	if err != nil {
		fmt.Fprintln(os.Stderr, bootstrap.ErrInit)
		return
	}
	server, err := bootstrap.NewMQTTServer(service)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if err := server.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
