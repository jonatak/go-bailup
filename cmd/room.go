package cmd

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/app"
)

type Room struct{}

func (r *Room) Run(ctx *app.AppContext) error {
	state, err := ctx.BailUp.GetState()
	if err != nil {
		return nil
	}

	for _, t := range state.Thermostats {
		fmt.Println(t.Name)
	}
	return nil
}
