package domain_test

import (
	"testing"

	"github.com/jonatak/go-bailup/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewThermostatRejectsInvalidPreset(t *testing.T) {
	heat := mustTemperatureSettings(t, domain.HVACSystemModeHeat, 20, 18)
	cool := mustTemperatureSettings(t, domain.HVACSystemModeCool, 24, 26)

	thermostat, err := domain.NewThermostat("Living Room", domain.ThermostatPreset("away"), true, false, heat, cool)

	require.ErrorIs(t, err, domain.ErrInvalidPresetMode)
	assert.Equal(t, domain.Thermostat{}, thermostat)
}

func TestNewThermostatRejectsInvalidTemperatureSettings(t *testing.T) {
	cool := mustTemperatureSettings(t, domain.HVACSystemModeCool, 24, 26)

	thermostat, err := domain.NewThermostat("Living Room", domain.PresetComfort, true, false, domain.TemperatureSettings{}, cool)

	require.ErrorIs(t, err, domain.ErrInvalidTemperatureRange)
	assert.Equal(t, domain.Thermostat{}, thermostat)
}

func TestThermostatValidateRejectsInvalidZeroValueThermostat(t *testing.T) {
	var thermostat domain.Thermostat

	require.ErrorIs(t, thermostat.Validate(), domain.ErrInvalidPresetMode)
}

func TestThermostatSetPresetValidatesPreset(t *testing.T) {
	thermostat := mustThermostat(t, "Living Room", domain.PresetComfort)

	err := thermostat.SetPreset(domain.PresetEco)

	require.NoError(t, err)
	assert.Equal(t, domain.PresetEco, thermostat.Preset())
}

func TestThermostatSetPresetDoesNotChangeStateForInvalidPreset(t *testing.T) {
	thermostat := mustThermostat(t, "Living Room", domain.PresetComfort)

	err := thermostat.SetPreset(domain.ThermostatPreset("away"))

	require.ErrorIs(t, err, domain.ErrInvalidPresetMode)
	assert.Equal(t, domain.PresetComfort, thermostat.Preset())
}

func mustTemperatureSettings(
	t *testing.T,
	mode domain.HVACSystemMode,
	comfort float64,
	eco float64,
) domain.TemperatureSettings {
	t.Helper()

	settings, err := domain.NewTemperatureSettings(mode, comfort, eco)
	require.NoError(t, err)

	return settings
}

func mustThermostat(t *testing.T, room string, preset domain.ThermostatPreset) domain.Thermostat {
	t.Helper()

	heat := mustTemperatureSettings(t, domain.HVACSystemModeHeat, 20, 18)
	cool := mustTemperatureSettings(t, domain.HVACSystemModeCool, 24, 26)
	thermostat, err := domain.NewThermostat(room, preset, true, false, heat, cool)
	require.NoError(t, err)

	return thermostat
}
