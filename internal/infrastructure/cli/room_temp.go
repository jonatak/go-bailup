package cli

import (
	"fmt"
	"strings"

	"github.com/jonatak/go-bailup/internal/app"
	"github.com/jonatak/go-bailup/internal/domain"
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
	system, err := appCtx.HVACService.CurrentState()
	if err != nil {
		return err
	}

	targetMode, err := targetHVACMode(system, mode)
	if err != nil {
		return err
	}

	targetPreset, err := targetThermostatPreset(system, roomName, preset)
	if err != nil {
		return err
	}

	if isDelta {
		current, err := system.Setpoint(roomName, targetMode, targetPreset)
		if err != nil {
			return err
		}
		value = current + value
	}

	if mode == "current" && preset == "current" {
		system, err = appCtx.HVACService.SetCurrentSetpoint(roomName, value)
	} else {
		system, err = appCtx.HVACService.SetTemperature(roomName, targetMode, targetPreset, value)
	}
	if err != nil {
		return err
	}

	thermostat, err := findDomainThermostat(system, roomName)
	if err != nil {
		return fmt.Errorf("%w after update", err)
	}

	fmt.Println("New Temperature setting:")
	fmt.Println(formatTemperatureSettings(thermostat))

	return nil
}

func targetHVACMode(system *domain.HVACSystem, mode string) (domain.HVACSystemMode, error) {
	if mode == "current" {
		return system.Mode(), nil
	}

	target := domain.HVACSystemMode(mode)
	if err := target.Validate(); err != nil {
		return "", err
	}

	return target, nil
}

func targetThermostatPreset(
	system *domain.HVACSystem,
	roomName string,
	preset string,
) (domain.ThermostatPreset, error) {
	if preset != "current" {
		target := domain.ThermostatPreset(preset)
		if err := target.Validate(); err != nil {
			return "", err
		}

		return target, nil
	}

	thermostat, err := findDomainThermostat(system, roomName)
	if err != nil {
		return "", err
	}

	return thermostat.Preset(), nil
}

func findDomainThermostat(system *domain.HVACSystem, roomName string) (domain.Thermostat, error) {
	for _, thermostat := range system.Thermostats() {
		if strings.EqualFold(thermostat.Room(), roomName) {
			return thermostat, nil
		}
	}

	return domain.Thermostat{}, fmt.Errorf("thermostat %q not found", roomName)
}
