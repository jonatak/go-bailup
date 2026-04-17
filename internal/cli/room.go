package cli

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/app"
)

type RoomTarget struct {
	Name string `arg:"" help:"Room name"`
}

type Room struct {
	On     RoomOn     `cmd:"" help:"Turn thermostat room on"`
	Off    RoomOff    `cmd:"" help:"Turn thermostat room off"`
	List   RoomList   `cmd:"" help:"List available room"`
	Preset RoomPreset `cmd:"" help:"Set thermostat preset in room"`
	Temp   RoomTemp   `cmd:"" help:"Manage room temperature"`
}

type RoomList struct{}

func (r *RoomList) Run(ctx *app.AppContext) error {
	state, err := ctx.BailUp.GetState()
	if err != nil {
		return err
	}

	for _, t := range state.Thermostats {
		fmt.Println(t.Name)
	}
	return nil
}
