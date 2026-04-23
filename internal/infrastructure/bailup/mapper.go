package bailup

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/application"
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

func CommandFromResolvedIntent(state *model.State, intent application.ResolvedIntent) (command.JSONCommand, error) {
	if state == nil {
		return nil, fmt.Errorf("state is nil")
	}
	if intent == nil {
		return nil, fmt.Errorf("intent is nil")
	}

	switch i := intent.(type) {
	case application.SetModeIntent:
		return command.NewHVACModeCommand(string(i.Mode))
	case application.SetRoomPresetIntent:
		return command.NewPresetCommand(state, i.Room, string(i.Preset))
	case application.SetRoomPowerIntent:
		return command.NewRoomPowerCommand(state, i.Room, i.On)
	case application.ResolvedSetTemperatureIntent:
		return temperatureCommandFromIntent(state, i)
	default:
		return nil, fmt.Errorf("unsupported resolved intent type %T", intent)
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

func temperatureCommandFromIntent(
	state *model.State,
	intent application.ResolvedSetTemperatureIntent,
) (command.JSONCommand, error) {
	thermostat := state.GetThermostatByName(intent.Room)
	if thermostat == nil {
		return nil, fmt.Errorf("map temperature intent for room %q: thermostat not found", intent.Room)
	}

	ucMode, err := model.UCModeFromString(string(intent.Mode))
	if err != nil {
		return nil, fmt.Errorf("map temperature intent for room %q: invalid HVAC mode %q: %w",
			intent.Room,
			intent.Mode,
			err,
		)
	}

	thMode, err := model.ThModeFromString(string(intent.Preset))
	if err != nil {
		return nil, fmt.Errorf("map temperature intent for room %q: invalid preset %q: %w",
			intent.Room,
			intent.Preset,
			err,
		)
	}

	return command.TemperatureCommand{
		ThermostatID: thermostat.ID,
		UCMode:       ucMode,
		ThMode:       thMode,
		Value:        intent.Value,
	}, nil
}
