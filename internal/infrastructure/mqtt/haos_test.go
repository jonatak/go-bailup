package mqtt

import (
	"encoding/json"
	"testing"

	"github.com/jonatak/go-bailup/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThermostatFromDomain(t *testing.T) {
	thermostat := testHAOSThermostat(t, 9152, "Salle Tv", domain.PresetComfort)

	got := ThermostatFromDomain(thermostat, "custom_bailup")

	assert.Equal(t, "Salle Tv Thermostat", got.Name)
	assert.Equal(t, "custom_bailup-salle_tv-9152", got.UniqueID)
	assert.Equal(t, "custom_bailup/th_9152/mode/set", got.ModeCommandTopic)
	assert.Equal(t, "custom_bailup/th_9152/mode", got.ModeStateTopic)
	assert.Equal(t, "custom_bailup/th_9152/temperature/set", got.TemperatureCommandTopic)
	assert.Equal(t, "custom_bailup/th_9152/temperature", got.TemperatureStateTopic)
	assert.Equal(t, "custom_bailup/th_9152/current_temperature", got.CurrentTemperatureTopic)
	assert.Equal(t, "custom_bailup/th_9152/preset_mode/set", got.PresetModeCommandTopic)
	assert.Equal(t, "custom_bailup/th_9152/preset_mode", got.PresetModeStateTopic)
	assert.Equal(t, "custom_bailup/th_9152/action", got.ActionTopic)
	assert.Equal(t, []string{"off", "auto"}, got.Modes)
	assert.Equal(t, []string{"eco", "comfort"}, got.PresetModes)
	assert.Equal(t, 16.0, got.MinTemp)
	assert.Equal(t, 30.0, got.MaxTemp)
	assert.Equal(t, 0.5, got.TempStep)
	assert.Equal(t, 0.1, got.Precision)
	assert.Equal(t, Device{
		Identifiers:   []string{"bailup_9152"},
		Manufacturer:  "Bail Industry",
		Name:          "Thermostat Salle Tv",
		SuggestedArea: "Salle Tv",
	}, got.Device)
}

func TestThermostatGeneralFromDomain(t *testing.T) {
	got := ThermostatGeneralFromDomain("custom_bailup")

	assert.Equal(t, "Thermostat General", got.Name)
	assert.Equal(t, "custom_bailup-general", got.UniqueID)
	assert.Equal(t, "custom_bailup/general/mode/set", got.ModeCommandTopic)
	assert.Equal(t, "custom_bailup/general/mode", got.ModeStateTopic)
	assert.Equal(t, "custom_bailup/general/current_temperature", got.CurrentTemperatureTopic)
	assert.Equal(t, []string{"off", "cool", "heat", "dry", "fan_only"}, got.Modes)
	assert.Equal(t, 0.1, got.Precision)
	assert.Equal(t, Device{
		Identifiers:  []string{"bailup_general"},
		Manufacturer: "Bail Industry",
		Name:         "Thermostat General",
	}, got.Device)
}

func TestThermostatJSONUsesSuggestedAreaTag(t *testing.T) {
	thermostat := testHAOSThermostat(t, 9152, "Salle Tv", domain.PresetComfort)

	payload, err := json.Marshal(ThermostatFromDomain(thermostat, "custom_bailup"))
	require.NoError(t, err)

	assert.Contains(t, string(payload), `"suggested_area":"Salle Tv"`)
	assert.NotContains(t, string(payload), `"suggested area"`)
}

func TestSlugify(t *testing.T) {
	assert.Equal(t, "salle_tv", slugify(" Salle Tv "))
}

func testHAOSThermostat(t *testing.T, id int, room string, preset domain.ThermostatPreset) domain.Thermostat {
	t.Helper()

	heatSetting, err := domain.NewTemperatureSettings(domain.HVACSystemModeHeat, 20, 18)
	require.NoError(t, err)
	coolSetting, err := domain.NewTemperatureSettings(domain.HVACSystemModeCool, 24, 26)
	require.NoError(t, err)

	thermostat, err := domain.NewThermostat(id, room, 20.0, preset, true, false, heatSetting, coolSetting)
	require.NoError(t, err)

	return thermostat
}
