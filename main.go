package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/jonatak/go-bailup/cmd"
	"github.com/jonatak/go-bailup/internal/app"
)

func main() {
	fmt.Println("Welcome to bailup.")

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

	ctx := kong.Parse(&cmd.CLI, kong.Bind(appCtx))

	err = ctx.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "an error occured: %v\n", err)
		return
	}

}
