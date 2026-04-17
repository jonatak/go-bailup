package cli

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/app"
	"github.com/jonatak/go-bailup/internal/bailup"
	"github.com/jonatak/go-bailup/internal/bailup/model"
)

type TemperatureTarget struct {
	Preset string `help:"Target preset" enum:"eco,comfort,current" default:"current"`
	Mode   string `help:"Target HVAC mode" enum:"heat,cool,current" default:"current"`
}

type TemperatureDelta struct {
	By float64 `help:"Temperature delta" default:"1"`
}

type RoomTemp struct {
	Set  RoomTempSet  `cmd:"" help:"Set room temperature"`
	Up   RoomTempUp   `cmd:"" help:"Increase room temperature"`
	Down RoomTempDown `cmd:"" help:"Decrease room temperature"`
}

type RoomTempSet struct {
	RoomTarget
	Value float64 `arg:"" help:"Target temperature"`
	TemperatureTarget
}

type RoomTempUp struct {
	RoomTarget
	TemperatureDelta
	TemperatureTarget
}

type RoomTempDown struct {
	RoomTarget
	TemperatureDelta
	TemperatureTarget
}

func (r *RoomTempSet) Run(appCtx *app.AppContext) error {
	return setRoomTemperature(appCtx, r.Name, r.Preset, r.Mode, r.Value, false)
}

func (r *RoomTempUp) Run(appCtx *app.AppContext) error {
	return setRoomTemperature(appCtx, r.Name, r.Preset, r.Mode, r.By, true)
}

func (r *RoomTempDown) Run(appCtx *app.AppContext) error {
	return setRoomTemperature(appCtx, r.Name, r.Preset, r.Mode, -r.By, true)
}

func setRoomTemperature(
	appCtx *app.AppContext,
	roomName string,
	preset string,
	mode string,
	value float64,
	isDelta bool,
) error {
	state, err := appCtx.BailUp.GetState()
	if err != nil {
		return err
	}

	th := state.GetThermostatByName(roomName)
	if th == nil {
		return fmt.Errorf("thermostat %q not found", roomName)
	}

	switch preset {
	case "eco":
		th.T1T2 = model.ThModeEco
	case "comfort":
		th.T1T2 = model.ThModeComfort
	}

	switch mode {
	case "current":
		if state.UCMode != model.UCModeCool && state.UCMode != model.UCModeHeat {
			return bailup.NewBailupError("UC is off, please specify temperature options --mode=heat or --mode=cool", nil)
		}
	case "heat":
		state.UCMode = model.UCModeHeat
	case "cool":
		state.UCMode = model.UCModeCool
	}

	if isDelta {
		current, err := currentSetpoint(state.UCMode, th)
		if err != nil {
			return err
		}
		value = current + value
	}

	cmd, err := bailup.NewTemperatureCommand(state, roomName, value)
	if err != nil {
		return err
	}

	state, err = appCtx.BailUp.Execute(cmd)
	if err != nil {
		return err
	}

	th = state.GetThermostatByName(roomName)
	if th == nil {
		return fmt.Errorf("thermostat %q not found after update", roomName)
	}

	fmt.Println("New Temperature setting:")
	fmt.Println(th.TemperatureSettingsString())

	return nil
}

func currentSetpoint(ucMode model.UCMode, th *model.Thermostat) (float64, error) {
	switch {
	case ucMode == model.UCModeHeat && th.T1T2 == model.ThModeComfort:
		return th.SetpointHotT1, nil
	case ucMode == model.UCModeHeat && th.T1T2 == model.ThModeEco:
		return th.SetpointHotT2, nil
	case ucMode == model.UCModeCool && th.T1T2 == model.ThModeComfort:
		return th.SetpointCoolT1, nil
	case ucMode == model.UCModeCool && th.T1T2 == model.ThModeEco:
		return th.SetpointCoolT2, nil
	default:
		return 0, fmt.Errorf("unsupported temperature combination: uc_mode=%s th_mode=%s", ucMode, th.T1T2)
	}
}
