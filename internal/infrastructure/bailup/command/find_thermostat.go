package command

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/model"
)

func findThermostat(s *model.State, room string) (*model.Thermostat, error) {
	thermostat := s.GetThermostatByName(room)
	if thermostat == nil {
		return nil, fmt.Errorf("thermostat %q not found", room)
	}
	return thermostat, nil
}
