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

	thermostat, err := domain.NewThermostat(1, "Living Room", 20.0, domain.ThermostatPreset("away"), true, false, heat, cool)

	require.ErrorIs(t, err, domain.ErrInvalidPresetMode)
	assert.Equal(t, domain.Thermostat{}, thermostat)
}

func TestNewThermostatRejectsInvalidTemperatureSettings(t *testing.T) {
	cool := mustTemperatureSettings(t, domain.HVACSystemModeCool, 24, 26)

	thermostat, err := domain.NewThermostat(1, "Living Room", 20.0, domain.PresetComfort, true, false, domain.TemperatureSettings{}, cool)

	require.ErrorIs(t, err, domain.ErrSetpointUnsupportedForMode)
	assert.Equal(t, domain.Thermostat{}, thermostat)
}

func TestThermostatValidateRejectsInvalidZeroValueThermostat(t *testing.T) {
	var thermostat domain.Thermostat

	require.ErrorIs(t, thermostat.Validate(), domain.ErrInvalidPresetMode)
}

func TestThermostatAction(t *testing.T) {
	heat := mustTemperatureSettings(t, domain.HVACSystemModeHeat, 20, 18)
	cool := mustTemperatureSettings(t, domain.HVACSystemModeCool, 24, 26)

	testCases := []struct {
		name      string
		isOn      bool
		isRunning bool
		mode      domain.HVACSystemMode
		want      domain.ThermostatAction
		wantErr   error
	}{
		{name: "off thermostat", isOn: false, isRunning: false, mode: domain.HVACSystemModeHeat, want: domain.ThermostatActionOff},
		{name: "idle thermostat", isOn: true, isRunning: false, mode: domain.HVACSystemModeHeat, want: domain.ThermostatActionIdle},
		{name: "cooling", isOn: true, isRunning: true, mode: domain.HVACSystemModeCool, want: domain.ThermostatActionCooling},
		{name: "heating", isOn: true, isRunning: true, mode: domain.HVACSystemModeHeat, want: domain.ThermostatActionHeating},
		{name: "drying", isOn: true, isRunning: true, mode: domain.HVACSystemModeDry, want: domain.ThermostatActionDrying},
		{name: "fan", isOn: true, isRunning: true, mode: domain.HVACSystemModeFanOnly, want: domain.ThermostatActionFan},
		{name: "running while off mode becomes idle", isOn: true, isRunning: true, mode: domain.HVACSystemModeOff, want: domain.ThermostatActionIdle},
		{name: "invalid mode", isOn: true, isRunning: true, mode: domain.HVACSystemMode("invalid"), wantErr: domain.ErrInvalidHVACMode},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			thermostat, err := domain.NewThermostat(1, "Living Room", 20.0, domain.PresetComfort, tc.isOn, tc.isRunning, heat, cool)
			require.NoError(t, err)

			got, err := thermostat.Action(tc.mode)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
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
	thermostat, err := domain.NewThermostat(1, room, 20.0, preset, true, false, heat, cool)
	require.NoError(t, err)

	return thermostat
}
