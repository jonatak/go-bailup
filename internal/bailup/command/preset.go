package command

import (
	"encoding/json"
	"fmt"
)

type Preset struct {
	ThermostatID int
	Value        int
}

func (p Preset) ToJSON() ([]byte, error) {
	payload := map[string]int{
		fmt.Sprintf("thermostats.%d.t1_t2", p.ThermostatID): p.Value,
	}
	return json.Marshal(payload)
}
