package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/jonatak/go-bailup/cmd"
	"github.com/jonatak/go-bailup/internal/app"
	kongcompletion "github.com/jotaen/kong-completion"
)

func main() {

	cliApp := kong.Must(&cmd.CLI{})
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

	appCtx, err := app.NewApp()

	if err != nil {
		if errors.Is(err, app.InitError) {
			fmt.Fprintln(os.Stderr, app.InitError)
		}
		return
	}

	if !appCtx.BailUp.IsConnected() {
		fmt.Println("Disconnected")
	}

	ctx.Bind(appCtx)

	err = ctx.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "an error occured: %v\n", err)
		return
	}

}
