package cmd

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/app"
	"github.com/jonatak/go-bailup/internal/bailup"
)

type RoomOn struct {
	RoomTarget
}

type RoomOff struct {
	RoomTarget
}

func (o *RoomOn) Run(appCtx *app.AppContext) error {
	return setRoomPower(appCtx, o.Name, true)
}

func (o *RoomOff) Run(appCtx *app.AppContext) error {
	return setRoomPower(appCtx, o.Name, false)
}

func setRoomPower(appCtx *app.AppContext, roomName string, on bool) error {
	state, err := appCtx.BailUp.GetState()
	if err != nil {
		return err
	}

	cmd, err := bailup.NewRoomPowerCommand(state, roomName, on)
	if err != nil {
		return err
	}

	state, err = appCtx.BailUp.Execute(cmd)
	if err != nil {
		return err
	}

	th := state.GetThermostatByName(roomName)
	status := "off"
	if th.IsOn {
		status = "on"
	}

	fmt.Printf("%s is now: %s\n", roomName, status)

	return nil
}
