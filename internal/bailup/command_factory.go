package bailup

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jonatak/go-bailup/internal/bailup/command"
	"github.com/jonatak/go-bailup/internal/bailup/model"
)

func NewCommand(s *model.State, room string, kind command.Kind, value string) (command.JSONCommand, error) {
	if kind == command.Mode && strings.EqualFold(room, "general") {
		ucMode, err := model.UCModeFromString(value)
		if err != nil {
			return nil, err
		}

		return command.ModeCommand{
			IsGeneral: true,
			Value:     int(ucMode),
		}, nil
	}

	var thermostat *model.Thermostat
	for i := range s.Thermostats {
		if strings.EqualFold(s.Thermostats[i].Name, room) {
			thermostat = &s.Thermostats[i]
			break
		}
	}
	if thermostat == nil {
		return nil, fmt.Errorf("thermostat %q not found", room)
	}

	switch kind {
	case command.PresetMode:
		thMode, err := model.ThModeFromString(value)
		if err != nil {
			return nil, err
		}
		return command.Preset{
			ThermostatID: thermostat.ID,
			Value:        int(thMode),
		}, nil

	case command.Mode:
		ucMode, err := model.UCModeFromString(value)
		if err != nil {
			return nil, err
		}
		return command.ModeCommand{
			ThermostatID: thermostat.ID,
			Value:        int(ucMode),
		}, nil

	case command.Temperature:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid temperature %q: %w", value, err)
		}
		return command.TemperatureCommand{
			ThermostatID: thermostat.ID,
			UCMode:       s.UCMode,
			ThMode:       thermostat.T1T2,
			Value:        floatValue,
		}, nil
	}

	return nil, fmt.Errorf("unsupported command %q", kind)
}
