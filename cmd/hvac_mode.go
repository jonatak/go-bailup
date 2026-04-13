package cmd

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/app"
	"github.com/jonatak/go-bailup/internal/bailup"
	"github.com/jonatak/go-bailup/internal/bailup/command"
)

type HvacMode struct {
	Mode string `arg:"" enum:"off,cool,heat,dry,fan-only"`
}

func (s *HvacMode) Run(appCtx *app.AppContext) error {
	state, err := appCtx.BailUp.GetState()
	if err != nil {
		return err
	}

	cmd, err := bailup.NewCommand(state, "general", command.Mode, s.Mode)
	if err != nil {
		return err
	}

	state, err = appCtx.BailUp.Execute(cmd)
	if err != nil {
		return err
	}

	fmt.Printf("HVAC mode is now: %s\n", state.UCMode.String())

	return nil
}
