package mqtt

import (
	"fmt"
	"strings"

	"github.com/jonatak/go-bailup/internal/domain"
)

type HAOSMode string

const (
	OFF      HAOSMode = "off"
	COOL     HAOSMode = "cool"
	HEAT     HAOSMode = "heat"
	DRY      HAOSMode = "dry"
	FAN_ONLY HAOSMode = "fan_only"
)

func PresetFromDomain(p domain.ThermostatPreset) string {
	return string(p)
}

func ModeFromDomain(m domain.HVACSystemMode) string {
	if m == domain.HVACSystemModeFanOnly {
		return "fan_only"
	}
	return string(m)
}

func ModeToDomain(m string) domain.HVACSystemMode {
	if m == "fan_only" {
		return domain.HVACSystemModeFanOnly
	}
	return domain.HVACSystemMode(m)
}

type Device struct {
	Identifiers   []string `json:"identifiers"`
	Manufacturer  string   `json:"manufacturer"`
	Name          string   `json:"name"`
	SuggestedArea string   `json:"suggested_area,omitempty"`
}

type MQTTThermostat struct {
	Name                    string `json:"name"`
	UniqueID                string `json:"unique_id"`
	ModeCommandTopic        string `json:"mode_command_topic"`
	ModeStateTopic          string `json:"mode_state_topic"`
	CurrentTemperatureTopic string `json:"current_temperature_topic"`
	TemperatureCommandTopic string `json:"temperature_command_topic,omitempty"`
	TemperatureStateTopic   string `json:"temperature_state_topic,omitempty"`
	PresetModeCommandTopic  string `json:"preset_mode_command_topic,omitempty"`
	PresetModeStateTopic    string `json:"preset_mode_state_topic,omitempty"`

	ActionTopic string   `json:"action_topic,omitempty"`
	Modes       []string `json:"modes"`
	PresetModes []string `json:"preset_modes,omitempty"`
	MinTemp     float64  `json:"min_temp,omitempty"`
	MaxTemp     float64  `json:"max_temp,omitempty"`
	TempStep    float64  `json:"temp_step,omitempty"`
	Precision   float64  `json:"precision"`
	Device      Device   `json:"device"`
}

func ThermostatFromDomain(t domain.Thermostat, prefix string) MQTTThermostat {
	return MQTTThermostat{
		Name:                    fmt.Sprintf("%s Thermostat", t.Room()),
		UniqueID:                fmt.Sprintf("%s-%s-%d", prefix, slugify(t.Room()), t.ID()),
		ModeCommandTopic:        fmt.Sprintf("%s/th_%d/mode/set", prefix, t.ID()),
		ModeStateTopic:          fmt.Sprintf("%s/th_%d/mode", prefix, t.ID()),
		TemperatureCommandTopic: fmt.Sprintf("%s/th_%d/temperature/set", prefix, t.ID()),
		TemperatureStateTopic:   fmt.Sprintf("%s/th_%d/temperature", prefix, t.ID()),
		CurrentTemperatureTopic: fmt.Sprintf("%s/th_%d/current_temperature", prefix, t.ID()),
		PresetModeCommandTopic:  fmt.Sprintf("%s/th_%d/preset_mode/set", prefix, t.ID()),
		PresetModeStateTopic:    fmt.Sprintf("%s/th_%d/preset_mode", prefix, t.ID()),
		ActionTopic:             fmt.Sprintf("%s/th_%d/action", prefix, t.ID()),
		Modes:                   []string{"off", "auto"},
		PresetModes:             []string{"eco", "comfort"},
		MinTemp:                 16,
		MaxTemp:                 30,
		TempStep:                0.5,
		Precision:               0.1,
		Device: Device{
			Identifiers:   []string{fmt.Sprintf("bailup_%d", t.ID())},
			Manufacturer:  "Bail Industry",
			Name:          fmt.Sprintf("Thermostat %s", t.Room()),
			SuggestedArea: t.Room(),
		},
	}
}

func ThermostatGeneralFromDomain(prefix string) MQTTThermostat {
	return MQTTThermostat{
		Name:                    "Thermostat General",
		UniqueID:                fmt.Sprintf("%s-general", prefix),
		ModeCommandTopic:        fmt.Sprintf("%s/general/mode/set", prefix),
		ModeStateTopic:          fmt.Sprintf("%s/general/mode", prefix),
		CurrentTemperatureTopic: fmt.Sprintf("%s/general/current_temperature", prefix),
		Modes:                   []string{"off", "cool", "heat", "dry", "fan_only"},
		Precision:               0.1,
		Device: Device{
			Identifiers:  []string{"bailup_general_unit"},
			Manufacturer: "Bail Industry",
			Name:         "Thermostat General",
		},
	}
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "_")
	return s
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.1f", f)
}
