package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/jonatak/go-bailup/internal/application"
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

func (r *RoomTempSet) Run(ctx context.Context, service *application.HVACService) error {
	return setRoomTemperature(ctx, service, r.Name, r.Preset, r.Mode, r.Value, false)
}

func (r *RoomTempUp) Run(ctx context.Context, service *application.HVACService) error {
	return setRoomTemperature(ctx, service, r.Name, r.Preset, r.Mode, r.By, true)
}

func (r *RoomTempDown) Run(ctx context.Context, service *application.HVACService) error {
	return setRoomTemperature(ctx, service, r.Name, r.Preset, r.Mode, -r.By, true)
}

func setRoomTemperature(
	ctx context.Context,
	service *application.HVACService,
	roomName string,
	preset string,
	mode string,
	value float64,
	isDelta bool,
) error {
	system, err := service.ApplyIntent(ctx, application.SetTemperatureIntent{
		Room:    roomName,
		Preset:  application.TemperaturePresetTarget(preset),
		Mode:    application.TemperatureModeTarget(mode),
		Value:   value,
		IsDelta: isDelta,
	})
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

func findDomainThermostat(system *domain.HVACSystem, roomName string) (domain.Thermostat, error) {
	for _, thermostat := range system.Thermostats() {
		if strings.EqualFold(thermostat.Room(), roomName) {
			return thermostat, nil
		}
	}

	return domain.Thermostat{}, fmt.Errorf("thermostat %q not found", roomName)
}
