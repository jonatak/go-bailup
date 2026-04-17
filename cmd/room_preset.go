package cmd

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/app"
	"github.com/jonatak/go-bailup/internal/bailup"
)

type RoomPreset struct {
	Preset string `arg:"" enum:"eco,comfort" help:"eco or comfort mode"`
	RoomTarget
}

func (r *RoomPreset) Run(appCtx *app.AppContext) error {
	state, err := appCtx.BailUp.GetState()
	if err != nil {
		return err
	}

	cmd, err := bailup.NewPresetCommand(state, r.Name, r.Preset)
	if err != nil {
		return err
	}

	state, err = appCtx.BailUp.Execute(cmd)
	if err != nil {
		return err
	}

	th := state.GetThermostatByName(r.Name)
	mode := th.T1T2.String()

	fmt.Printf("%s mode is now: %s\n", r.Name, mode)

	return nil
}
