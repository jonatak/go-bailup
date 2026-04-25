package cli

import (
	"testing"

	"github.com/jonatak/go-bailup/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatHVACSystem(t *testing.T) {
	system := testFormatHVACSystem(t)

	got := formatHVACSystem(system)

	assert.Contains(t, got, "Unit mode: heat")
	assert.Contains(t, got, "Thermostats (2):")
	assert.Contains(t, got, "#1 Living Room")
	assert.Contains(t, got, "Temperature: 20.0 C")
	assert.Contains(t, got, "Status: "+highlight("on", ansiGreen)+", running=false")
	assert.Contains(t, got, "Active preset: "+highlight("comfort", ansiGreen))
	assert.Contains(t, got, "Heat setpoints: comfort="+highlight("20.0 C", ansiCyan)+", eco=18.0 C")
	assert.Contains(t, got, "Cool setpoints: comfort="+highlight("24.0 C", ansiCyan)+", eco=26.0 C")
	assert.Contains(t, got, "#2 Bedroom")
	assert.Contains(t, got, "Temperature: 20.0 C")
	assert.Contains(t, got, "Status: off, running=true")
	assert.Contains(t, got, "Active preset: "+highlight("eco", ansiGreen))
}

func TestFormatHVACSystemWithoutThermostats(t *testing.T) {
	system, err := domain.NewHVACSystem(domain.HVACSystemModeOff, nil)
	require.NoError(t, err)

	got := formatHVACSystem(system)

	assert.Equal(t, "Unit mode: off\nThermostats: none", got)
}

func TestFormatHVACSystemNil(t *testing.T) {
	assert.Equal(t, "HVAC system: unavailable", formatHVACSystem(nil))
}

func testFormatHVACSystem(t *testing.T) *domain.HVACSystem {
	t.Helper()

	livingRoom := testFormatThermostat(t, "Living Room", domain.PresetComfort, true, false)
	bedroom := testFormatThermostat(t, "Bedroom", domain.PresetEco, false, true)
	system, err := domain.NewHVACSystem(domain.HVACSystemModeHeat, []domain.Thermostat{
		livingRoom,
		bedroom,
	})
	require.NoError(t, err)

	return system
}

func testFormatThermostat(
	t *testing.T,
	room string,
	preset domain.ThermostatPreset,
	isOn bool,
	isRunning bool,
) domain.Thermostat {
	t.Helper()

	heatSetting, err := domain.NewTemperatureSettings(domain.HVACSystemModeHeat, 20, 18)
	require.NoError(t, err)
	coolSetting, err := domain.NewTemperatureSettings(domain.HVACSystemModeCool, 24, 26)
	require.NoError(t, err)
	thermostat, err := domain.NewThermostat(
		1,
		room,
		20.0,
		preset,
		isOn,
		isRunning,
		heatSetting,
		coolSetting,
	)
	require.NoError(t, err)

	return thermostat
}
