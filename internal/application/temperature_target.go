package application

import (
	"fmt"
	"strings"

	"github.com/jonatak/go-bailup/internal/domain"
)

type resolvedTemperatureTarget struct {
	room   string
	mode   domain.HVACSystemMode
	preset domain.ThermostatPreset
	value  float64
}

func resolveTemperatureTarget(
	system *domain.HVACSystem,
	intent SetTemperatureIntent,
) (resolvedTemperatureTarget, error) {
	targetMode, err := targetHVACMode(system, intent.Mode)
	if err != nil {
		return resolvedTemperatureTarget{}, err
	}

	targetPreset, err := targetThermostatPreset(system, intent.Room, intent.Preset)
	if err != nil {
		return resolvedTemperatureTarget{}, err
	}

	value := intent.Value
	if intent.IsDelta {
		current, err := system.Setpoint(intent.Room, targetMode, targetPreset)
		if err != nil {
			return resolvedTemperatureTarget{}, err
		}
		value = current + value
	}

	return resolvedTemperatureTarget{
		room:   intent.Room,
		mode:   targetMode,
		preset: targetPreset,
		value:  value,
	}, nil
}

func (r resolvedTemperatureTarget) Intent() ResolvedSetTemperatureIntent {
	return ResolvedSetTemperatureIntent{
		Room:   r.room,
		Preset: r.preset,
		Mode:   r.mode,
		Value:  r.value,
	}
}

func targetHVACMode(system *domain.HVACSystem, mode TemperatureModeTarget) (domain.HVACSystemMode, error) {
	if mode == TemperatureModeCurrent {
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
	preset TemperaturePresetTarget,
) (domain.ThermostatPreset, error) {
	if preset != TemperaturePresetCurrent {
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
