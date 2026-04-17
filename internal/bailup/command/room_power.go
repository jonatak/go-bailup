package command

import (
	"encoding/json"
	"fmt"
)

type RoomPowerCommand struct {
	ThermostatID int
	On           bool
}

func (r RoomPowerCommand) ToJSON() ([]byte, error) {
	payload := map[string]bool{
		fmt.Sprintf("thermostats.%d.is_on", r.ThermostatID): r.On,
	}
	return json.Marshal(payload)
}
