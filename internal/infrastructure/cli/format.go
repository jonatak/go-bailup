package cli

import (
	"fmt"
	"strings"

	"github.com/jonatak/go-bailup/internal/domain"
)

const (
	ansiReset = "\033[0m"
	ansiBold  = "\033[1m"
	ansiGreen = "\033[32m"
	ansiCyan  = "\033[36m"
)

func formatHVACSystem(system *domain.HVACSystem) string {
	if system == nil {
		return "HVAC system: unavailable"
	}

	var b strings.Builder

	fmt.Fprintf(&b, "Unit mode: %s\n", system.Mode())

	thermostats := system.Thermostats()
	if len(thermostats) == 0 {
		b.WriteString("Thermostats: none")
		return b.String()
	}

	fmt.Fprintf(&b, "Thermostats (%d):\n", len(thermostats))
	for i, thermostat := range thermostats {
		fmt.Fprintf(&b, "- %s\n", strings.ReplaceAll(formatThermostat(i+1, thermostat), "\n", "\n  "))
	}

	return strings.TrimRight(b.String(), "\n")
}

func formatThermostat(number int, thermostat domain.Thermostat) string {
	var b strings.Builder

	fmt.Fprintf(&b, "#%d %s\n", number, thermostat.Room())
	fmt.Fprintf(&b, "  Status: %s, running=%t\n", formatPowerStatus(thermostat.IsOn()), thermostat.IsRunning())
	fmt.Fprintf(&b, "  Active preset: %s\n", highlight(string(thermostat.Preset()), ansiGreen))
	fmt.Fprintf(&b, "%s", formatTemperatureSettings(thermostat))

	return strings.TrimRight(b.String(), "\n")
}

func formatTemperatureSettings(thermostat domain.Thermostat) string {
	var b strings.Builder

	heatSetting := thermostat.HeatSetting()
	coolSetting := thermostat.CoolSetting()

	fmt.Fprintf(&b, "  Heat setpoints: comfort=%s, eco=%s\n",
		formatSetpoint(domain.PresetComfort, thermostat.Preset(), heatSetting.Comfort()),
		formatSetpoint(domain.PresetEco, thermostat.Preset(), heatSetting.Eco()),
	)
	fmt.Fprintf(&b, "  Cool setpoints: comfort=%s, eco=%s\n",
		formatSetpoint(domain.PresetComfort, thermostat.Preset(), coolSetting.Comfort()),
		formatSetpoint(domain.PresetEco, thermostat.Preset(), coolSetting.Eco()),
	)

	return strings.TrimRight(b.String(), "\n")
}

func formatSetpoint(preset domain.ThermostatPreset, activePreset domain.ThermostatPreset, value float64) string {
	setpoint := fmt.Sprintf("%.1f C", value)
	if activePreset == preset {
		return highlight(setpoint, ansiCyan)
	}
	return setpoint
}

func formatPowerStatus(isOn bool) string {
	if isOn {
		return highlight("on", ansiGreen)
	}
	return "off"
}

func highlight(value string, color string) string {
	return color + ansiBold + value + ansiReset
}
