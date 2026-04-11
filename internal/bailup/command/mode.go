package command

import (
	"encoding/json"
	"fmt"
)

type ModeCommand struct {
	ThermostatID int
	IsGeneral    bool
	Value        int
}

func (m ModeCommand) ToJSON() ([]byte, error) {
	payload := map[string]any{}

	if m.IsGeneral {
		payload["uc_mode"] = m.Value
		return json.Marshal(payload)
	}

	switch m.Value {
	case 0:
		payload[fmt.Sprintf("thermostats.%d.is_on", m.ThermostatID)] = false
	case 5:
		payload[fmt.Sprintf("thermostats.%d.is_on", m.ThermostatID)] = true
	default:
		return nil, fmt.Errorf("unsupported thermostat mode command %d", m.Value)
	}

	return json.Marshal(payload)
}
