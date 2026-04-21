package bailup

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/domain"
	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/command"
	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/model"
)

func HVACSystemFromState(state *model.State) (*domain.HVACSystem, error) {
	if state == nil {
		return nil, fmt.Errorf("state is nil")
	}

	thermostats := make([]domain.Thermostat, 0, len(state.Thermostats))
	for _, thermostat := range state.Thermostats {
		domainThermostat, err := thermostatFromModel(thermostat)
		if err != nil {
			return nil, err
		}
		thermostats = append(thermostats, domainThermostat)
	}

	return domain.NewHVACSystem(
		domain.HVACSystemMode(state.UCMode.String()),
		thermostats,
	)
}

func CommandFromChange(state *model.State, change domain.Change) (command.JSONCommand, error) {
	if state == nil {
		return nil, fmt.Errorf("state is nil")
	}
	if change == nil {
		return nil, fmt.Errorf("change is nil")
	}

	switch c := change.(type) {
	case domain.HVACModeChanged:
		return command.NewHVACModeCommand(string(c.Mode))
	case domain.RoomPresetChanged:
		return command.NewPresetCommand(state, c.Room, string(c.Preset))
	case domain.RoomPowerChanged:
		return command.NewRoomPowerCommand(state, c.Room, c.On)
	case domain.TemperatureChanged:
		return temperatureCommandFromChange(state, c)
	default:
		return nil, fmt.Errorf("unsupported domain change kind %q", change.Kind())
	}
}

func thermostatFromModel(thermostat model.Thermostat) (domain.Thermostat, error) {
	heatSetting, err := domain.NewTemperatureSettings(
		domain.HVACSystemModeHeat,
		thermostat.SetpointHotT1,
		thermostat.SetpointHotT2,
	)
	if err != nil {
		return domain.Thermostat{}, fmt.Errorf(
			"map thermostat %q heat settings: comfort=%.1f eco=%.1f: %w",
			thermostat.Name,
			thermostat.SetpointHotT1,
			thermostat.SetpointHotT2,
			err,
		)
	}

	coolSetting, err := domain.NewTemperatureSettings(
		domain.HVACSystemModeCool,
		thermostat.SetpointCoolT1,
		thermostat.SetpointCoolT2,
	)
	if err != nil {
		return domain.Thermostat{}, fmt.Errorf(
			"map thermostat %q cool settings: comfort=%.1f eco=%.1f: %w",
			thermostat.Name,
			thermostat.SetpointCoolT1,
			thermostat.SetpointCoolT2,
			err,
		)
	}

	return domain.NewThermostat(
		thermostat.Name,
		domain.ThermostatPreset(thermostat.T1T2.String()),
		thermostat.IsOn,
		thermostat.MotorState > 4,
		heatSetting,
		coolSetting,
	)
}

func temperatureCommandFromChange(
	state *model.State,
	change domain.TemperatureChanged,
) (command.JSONCommand, error) {
	thermostat := state.GetThermostatByName(change.Room)
	if thermostat == nil {
		return nil, fmt.Errorf("map temperature change for room %q: thermostat not found", change.Room)
	}

	ucMode, err := model.UCModeFromString(string(change.Mode))
	if err != nil {
		return nil, fmt.Errorf("map temperature change for room %q: invalid HVAC mode %q: %w",
			change.Room,
			change.Mode,
			err,
		)
	}

	thMode, err := model.ThModeFromString(string(change.Preset))
	if err != nil {
		return nil, fmt.Errorf("map temperature change for room %q: invalid preset %q: %w",
			change.Room,
			change.Preset,
			err,
		)
	}

	return command.TemperatureCommand{
		ThermostatID: thermostat.ID,
		UCMode:       ucMode,
		ThMode:       thMode,
		Value:        change.Value,
	}, nil
}
