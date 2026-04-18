package command_test

import (
	"encoding/json"
	"testing"

	"github.com/jonatak/go-bailup/internal/bailup/command"
	"github.com/jonatak/go-bailup/internal/bailup/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemperatureCommandToJSONHeatComfort(t *testing.T) {
	cmd := command.TemperatureCommand{
		ThermostatID: 42,
		Value:        20.5,
		UCMode:       model.UCModeHeat,
		ThMode:       model.ThModeComfort,
	}

	assertTemperaturePayload(t, cmd, map[string]float64{
		"thermostats.42.setpoint_hot_t1": 20.5,
	})
}

func TestTemperatureCommandToJSONHeatEco(t *testing.T) {
	cmd := command.TemperatureCommand{
		ThermostatID: 42,
		Value:        20.5,
		UCMode:       model.UCModeHeat,
		ThMode:       model.ThModeEco,
	}

	assertTemperaturePayload(t, cmd, map[string]float64{
		"thermostats.42.setpoint_hot_t2": 20.5,
	})
}

func TestTemperatureCommandToJSONCoolComfort(t *testing.T) {
	cmd := command.TemperatureCommand{
		ThermostatID: 42,
		Value:        20.5,
		UCMode:       model.UCModeCool,
		ThMode:       model.ThModeComfort,
	}

	assertTemperaturePayload(t, cmd, map[string]float64{
		"thermostats.42.setpoint_cool_t1": 20.5,
	})
}

func TestTemperatureCommandToJSONCoolEco(t *testing.T) {
	cmd := command.TemperatureCommand{
		ThermostatID: 42,
		Value:        20.5,
		UCMode:       model.UCModeCool,
		ThMode:       model.ThModeEco,
	}

	assertTemperaturePayload(t, cmd, map[string]float64{
		"thermostats.42.setpoint_cool_t2": 20.5,
	})
}

func TestTemperatureCommandToJSONRoundsValueToOneDecimal(t *testing.T) {
	cmd := command.TemperatureCommand{
		ThermostatID: 42,
		Value:        20.56,
		UCMode:       model.UCModeHeat,
		ThMode:       model.ThModeComfort,
	}

	assertTemperaturePayload(t, cmd, map[string]float64{
		"thermostats.42.setpoint_hot_t1": 20.6,
	})
}

func TestTemperatureCommandToJSONReturnsErrorForUnsupportedModeCombination(t *testing.T) {
	cmd := command.TemperatureCommand{
		ThermostatID: 42,
		Value:        20.5,
		UCMode:       model.UCModeOff,
		ThMode:       model.ThModeComfort,
	}

	got, err := cmd.ToJSON()

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "unsupported temperature combination")
}

func assertTemperaturePayload(
	t *testing.T,
	cmd command.TemperatureCommand,
	want map[string]float64,
) {
	t.Helper()

	got, err := cmd.ToJSON()
	require.NoError(t, err)

	var payload map[string]float64
	require.NoError(t, json.Unmarshal(got, &payload))

	assert.Equal(t, want, payload)
}
