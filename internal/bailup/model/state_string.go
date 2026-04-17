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
		fmt.Fprintf(&b, "- %s\n", strings.ReplaceAll(t.String(), "\n", "\n  "))
	}

	return strings.TrimRight(b.String(), "\n")
}
