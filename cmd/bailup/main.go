package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/jonatak/go-bailup/internal/bootstrap"
	"github.com/jonatak/go-bailup/internal/handler/cli"
	kongcompletion "github.com/jotaen/kong-completion"
)

func main() {

	cliApp := kong.Must(&cli.CLI{})
	kongcompletion.Register(cliApp)

	ctx, err := cliApp.Parse(os.Args[1:])
	if err != nil {
		cliApp.Printf("%s", err)
		cliApp.Exit(1)
		return
	}

	if ctx.Command() == "completion" {
		if err := ctx.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "an error occured: %v\n", err)
		}
		return
	}

	if ctx.Command() != "serve" {
		service, err := bootstrap.NewHVACService()

		if err != nil {
			if errors.Is(err, bootstrap.InitError) {
				fmt.Fprintln(os.Stderr, bootstrap.InitError)
			}
			return
		}
		ctx.Bind(service)
	}

	err = ctx.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "an error occured: %v\n", err)
		return
	}

}
