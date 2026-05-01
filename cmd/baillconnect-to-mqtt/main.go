package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonatak/baillconnect-to-mqtt/internal/bootstrap"
	"github.com/jonatak/baillconnect-to-mqtt/internal/config"
)

func main() {
	configPath := flag.String("config", "", "path to configuration file")
	flag.Parse()

	slog.Info(fmt.Sprintf("Start baillconnect-to-mqtt version:%s, commit:%s, buildtime:%s", config.Version, config.CommitSHA, config.BuildTime))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	service, err := bootstrap.NewHVACService(cfg)

	if err != nil {
		fmt.Fprintln(os.Stderr, bootstrap.ErrInit)
		return
	}
	server, err := bootstrap.NewMQTTServer(service, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if err := server.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
