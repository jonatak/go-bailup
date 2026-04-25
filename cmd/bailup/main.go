package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/jonatak/go-bailup/internal/application"
	"github.com/jonatak/go-bailup/internal/bootstrap"
	"github.com/jonatak/go-bailup/internal/infrastructure/cli"
	kongcompletion "github.com/jotaen/kong-completion"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	cliApp := kong.Must(&cli.CLI{}, kong.BindTo(ctx, (*context.Context)(nil)))
	kongcompletion.Register(cliApp)

	kongCtx, err := cliApp.Parse(os.Args[1:])
	if err != nil {
		cliApp.Printf("%s", err)
		cliApp.Exit(1)
		return
	}

	if kongCtx.Command() == "completion" {
		if err := kongCtx.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "an error occured: %v\n", err)
		}
		return
	}
	service, err := bootstrap.NewHVACService()

	if err != nil {
		fmt.Fprintln(os.Stderr, bootstrap.InitError)
		return
	}
	if kongCtx.Command() == "serve" {
		server, err := bootstrap.NewMQTTServer(service)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		kongCtx.BindTo(server, (*application.Server)(nil))
	} else {
		kongCtx.Bind(service)
	}

	err = kongCtx.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "an error occured: %v\n", err)
		return
	}

}
