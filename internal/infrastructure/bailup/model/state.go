package model

import (
	"strings"
)

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

func (s *State) GetThermostatByName(name string) *Thermostat {
	var thermostat *Thermostat
	for i := range s.Thermostats {
		if strings.EqualFold(s.Thermostats[i].Name, name) {
			thermostat = &s.Thermostats[i]
			break
		}
	}
	return thermostat
}
