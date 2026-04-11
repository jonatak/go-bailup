package bailup

import (
	"fmt"
	"strings"
)

type UCMode int

const (
	UCModeOff UCMode = iota
	UCModeCool
	UCModeHeat
	UCModeDry
	UCModeFanOnly
	UCModeAuto
)

type ThMode int

const (
	ThModeComfort ThMode = 1 + iota
	ThModeEco
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

type State struct {
	ID          int          `json:"id"`
	Mbus        bool         `json:"mbus"`
	UCMode      UCMode       `json:"uc_mode"`
	UCHotMin    int          `json:"uc_hot_min"`
	UCHotMax    int          `json:"uc_hot_max"`
	UCColdMin   int          `json:"uc_cold_min"`
	UCColdMax   int          `json:"uc_cold_max"`
	IsConnected bool         `json:"is_connected"`
	Thermostats []Thermostat `json:"thermostats"`
}

func (m UCMode) String() string {
	switch m {
	case UCModeOff:
		return "off"
	case UCModeCool:
		return "cool"
	case UCModeHeat:
		return "heat"
	case UCModeDry:
		return "dry"
	case UCModeFanOnly:
		return "fan-only"
	case UCModeAuto:
		return "auto"
	default:
		return fmt.Sprintf("unknown(%d)", int(m))
	}
}

func (m ThMode) String() string {
	switch m {
	case ThModeComfort:
		return "comfort"
	case ThModeEco:
		return "eco"
	default:
		return fmt.Sprintf("unknown(%d)", int(m))
	}
}

func (s State) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "State #%d\n", s.ID)
	fmt.Fprintf(&b, "Connected: %t\n", s.IsConnected)
	fmt.Fprintf(&b, "MBus: %t\n", s.Mbus)
	fmt.Fprintf(&b, "Unit mode: %s\n", s.UCMode)
	fmt.Fprintf(&b, "Heat range: %d-%d C\n", s.UCHotMin, s.UCHotMax)
	fmt.Fprintf(&b, "Cold range: %d-%d C\n", s.UCColdMin, s.UCColdMax)

	if len(s.Thermostats) == 0 {
		b.WriteString("Thermostats: none")
		return b.String()
	}

	fmt.Fprintf(&b, "Thermostats (%d):\n", len(s.Thermostats))
	for _, t := range s.Thermostats {
		fmt.Fprintf(&b, "- #%d %s\n", t.Number, t.Name)
		fmt.Fprintf(&b, "  id=%d key=%s zone=%d\n", t.ID, t.Key, t.Zone)
		fmt.Fprintf(&b, "  connected=%t on=%t battery_low=%t\n", t.IsConnected, t.IsOn, t.IsBatteryLow)
		fmt.Fprintf(&b, "  temperature=%.1f C mode=%s motor_state=%d\n", t.Temperature, t.T1T2, t.MotorState)
		fmt.Fprintf(&b, "  hot_setpoints=%.1f / %.1f C\n", t.SetpointHotT1, t.SetpointHotT2)
		fmt.Fprintf(&b, "  cool_setpoints=%.1f / %.1f C\n", t.SetpointCoolT1, t.SetpointCoolT2)
	}

	return strings.TrimRight(b.String(), "\n")
}
