package model

import (
	"fmt"
	"strings"
)

const (
	ansiReset = "\033[0m"
	ansiBold  = "\033[1m"
	ansiGreen = "\033[32m"
	ansiCyan  = "\033[36m"
)

type Response struct {
	Data State `json:"data"`
}

type Thermostat struct {
	ID             int     `json:"id"`
	Key            string  `json:"key"`
	Number         int     `json:"number"`
	Name           string  `json:"name"`
	Temperature    float64 `json:"temperature"`
	Zone           int     `json:"zone"`
	IsOn           bool    `json:"is_on"`
	SetpointHotT1  float64 `json:"setpoint_hot_t1"`
	SetpointHotT2  float64 `json:"setpoint_hot_t2"`
	SetpointCoolT1 float64 `json:"setpoint_cool_t1"`
	SetpointCoolT2 float64 `json:"setpoint_cool_t2"`
	MotorState     int     `json:"motor_state"`
	T1T2           ThMode  `json:"t1_t2"`
	IsBatteryLow   bool    `json:"is_battery_low"`
	IsConnected    bool    `json:"is_connected"`
}

func (t Thermostat) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "#%d %s\n", t.Number, t.Name)
	fmt.Fprintf(&b, "  Status: %s, connected=%t, battery_low=%t\n", t.powerStatus(), t.IsConnected, t.IsBatteryLow)
	fmt.Fprintf(&b, "  Current temperature: %.1f C\n", t.Temperature)
	fmt.Fprintf(&b, "  Active preset: %s\n", highlight(t.T1T2.String(), ansiGreen))
	fmt.Fprintf(&b, "  Zone: %d, motor_state=%d\n", t.Zone, t.MotorState)
	fmt.Fprintf(&b, "%s", t.TemperatureSettingsString())

	return strings.TrimRight(b.String(), "\n")
}

func (t Thermostat) TemperatureSettingsString() string {
	var b strings.Builder

	fmt.Fprintf(&b, "  Heat setpoints: comfort=%s, eco=%s\n",
		t.formatSetpoint(ThModeComfort, t.SetpointHotT1),
		t.formatSetpoint(ThModeEco, t.SetpointHotT2),
	)
	fmt.Fprintf(&b, "  Cool setpoints: comfort=%s, eco=%s\n",
		t.formatSetpoint(ThModeComfort, t.SetpointCoolT1),
		t.formatSetpoint(ThModeEco, t.SetpointCoolT2),
	)

	return strings.TrimRight(b.String(), "\n")
}

func (t Thermostat) formatSetpoint(mode ThMode, value float64) string {
	setpoint := fmt.Sprintf("%.1f C", value)
	if t.T1T2 == mode {
		return highlight(setpoint, ansiCyan)
	}
	return setpoint
}

func (t Thermostat) powerStatus() string {
	if t.IsOn {
		return highlight("on", ansiGreen)
	}
	return "off"
}

func highlight(value string, color string) string {
	return color + ansiBold + value + ansiReset
}
