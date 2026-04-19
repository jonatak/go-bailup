package command

import (
	"encoding/json"
	"fmt"

	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/model"
)

type Preset struct {
	ThermostatID int
	Value        int
}

func NewPresetCommand(s *model.State, room string, preset string) (JSONCommand, error) {
	thermostat, err := findThermostat(s, room)
	if err != nil {
		return nil, err
	}

	thMode, err := model.ThModeFromString(preset)
	if err != nil {
		return nil, err
	}

	return Preset{
		ThermostatID: thermostat.ID,
		Value:        int(thMode),
	}, nil
}

func (p Preset) ToJSON() ([]byte, error) {
	payload := map[string]int{
		fmt.Sprintf("thermostats.%d.t1_t2", p.ThermostatID): p.Value,
	}
	return json.Marshal(payload)
}
