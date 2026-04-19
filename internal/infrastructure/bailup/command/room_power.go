package command

import (
	"encoding/json"
	"fmt"

	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/model"
)

type RoomPowerCommand struct {
	ThermostatID int
	On           bool
}

func NewRoomPowerCommand(s *model.State, room string, on bool) (JSONCommand, error) {
	thermostat, err := findThermostat(s, room)
	if err != nil {
		return nil, err
	}

	return RoomPowerCommand{
		ThermostatID: thermostat.ID,
		On:           on,
	}, nil
}

func (r RoomPowerCommand) ToJSON() ([]byte, error) {
	payload := map[string]bool{
		fmt.Sprintf("thermostats.%d.is_on", r.ThermostatID): r.On,
	}
	return json.Marshal(payload)
}
