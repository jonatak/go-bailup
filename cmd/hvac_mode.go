package cmd

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/app"
	"github.com/jonatak/go-bailup/internal/bailup"
)

type HvacMode struct {
	Mode string `arg:"" enum:"off,cool,heat,dry,fan-only"`
}

func (s *HvacMode) Run(appCtx *app.AppContext) error {
	cmd, err := bailup.NewHVACModeCommand(s.Mode)
	if err != nil {
		return err
	}

	state, err := appCtx.BailUp.Execute(cmd)
	if err != nil {
		return err
	}

	fmt.Printf("HVAC mode is now: %s\n", state.UCMode.String())

	return nil
}
