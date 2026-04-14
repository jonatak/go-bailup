package cmd

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/app"
)

type Room struct {
	On     RoomOn     `cmd:"" help:"Turn thermostat room on"`
	Off    RoomOff    `cmd:"" help:"Turm thermostat room off"`
	List   RoomList   `cmd:"" help:"List available room"`
	Preset RoomPreset `cmd:"" help:"Set thermostat preset in room"`
	Temp   RoomTemp   `cmd:"" help:"Manage room temperature"`
}

type RoomOn struct {
	Name string `arg:"" help:"Room name"`
}

type RoomOff struct {
	Name string `arg:"" help:"Room name"`
}

type RoomList struct {
}

type RoomPreset struct {
	Preset string `arg:"" enum:"eco,comfort" help:"eco or comfort mode"`
	Name   string `arg:"" help:"Room name"`
}

type RoomTemp struct {
	Set  RoomTempSet  `cmd:"" help:"Set room temperature"`
	Up   RoomTempUp   `cmd:"" help:"Increase room temperature"`
	Down RoomTempDown `cmd:"" help:"Decrease room temperature"`
}

type RoomTempSet struct {
	Name   string  `arg:"" help:"Room name"`
	Value  float64 `arg:"" help:"Target temperature"`
	Preset string  `help:"Target preset" enum:"eco,comfort,current" default:"current"`
}

type RoomTempUp struct {
	Name   string  `arg:"" help:"Room name"`
	By     float64 `help:"Temperature delta" default:"1"`
	Preset string  `help:"Target preset" enum:"eco,comfort,current" default:"current"`
}

type RoomTempDown struct {
	Name   string  `arg:"" help:"Room name"`
	By     float64 `help:"Temperature delta" default:"1"`
	Preset string  `help:"Target preset" enum:"eco,comfort,current" default:"current"`
}

func (r *RoomList) Run(ctx *app.AppContext) error {
	state, err := ctx.BailUp.GetState()
	if err != nil {
		return nil
	}

	for _, t := range state.Thermostats {
		fmt.Println(t.Name)
	}
	return nil
}

func (o *RoomOn) Run(ctx *app.AppContext) error {
	fmt.Println("Room on")
	return nil
}
