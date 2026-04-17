package bailup

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/bailup/command"
	"github.com/jonatak/go-bailup/internal/bailup/model"
)

func NewHVACModeCommand(mode string) (command.JSONCommand, error) {
	ucMode, err := model.UCModeFromString(mode)
	if err != nil {
		return nil, err
	}

	return command.ModeCommand{
		Value: int(ucMode),
	}, nil
}

func NewRoomPowerCommand(s *model.State, room string, on bool) (command.JSONCommand, error) {
	thermostat, err := findThermostat(s, room)
	if err != nil {
		return nil, err
	}

	return command.RoomPowerCommand{
		ThermostatID: thermostat.ID,
		On:           on,
	}, nil
}

func NewPresetCommand(s *model.State, room string, preset string) (command.JSONCommand, error) {
	thermostat, err := findThermostat(s, room)
	if err != nil {
		return nil, err
	}

	thMode, err := model.ThModeFromString(preset)
	if err != nil {
		return nil, err
	}

	return command.Preset{
		ThermostatID: thermostat.ID,
		Value:        int(thMode),
	}, nil
}

func NewTemperatureCommand(s *model.State, room string, value float64) (command.JSONCommand, error) {
	thermostat, err := findThermostat(s, room)
	if err != nil {
		return nil, err
	}

	return command.TemperatureCommand{
		ThermostatID: thermostat.ID,
		UCMode:       s.UCMode,
		ThMode:       thermostat.T1T2,
		Value:        value,
	}, nil
}

func findThermostat(s *model.State, room string) (*model.Thermostat, error) {
	thermostat := s.GetThermostatByName(room)
	if thermostat == nil {
		return nil, fmt.Errorf("thermostat %q not found", room)
	}
	return thermostat, nil
}
