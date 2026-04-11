package model

import (
	"fmt"
	"strings"
)

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
