package domain

import "strings"

type HVACSystem struct {
	mode        HVACSystemMode
	thermostats []Thermostat
}

func NewHVACSystem(mode HVACSystemMode, thermostats []Thermostat) (*HVACSystem, error) {
	system := &HVACSystem{
		mode:        mode,
		thermostats: append([]Thermostat(nil), thermostats...),
	}

	if err := system.Validate(); err != nil {
		return nil, err
	}

	return system, nil
}

func (s *HVACSystem) Validate() error {
	if err := s.mode.Validate(); err != nil {
		return err
	}

	for i := range s.thermostats {
		if err := s.thermostats[i].Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (s *HVACSystem) Mode() HVACSystemMode {
	return s.mode
}

func (s *HVACSystem) Thermostats() []Thermostat {
	return append([]Thermostat(nil), s.thermostats...)
}

func (s *HVACSystem) SetMode(mode HVACSystemMode) (Change, error) {
	if err := mode.Validate(); err != nil {
		return nil, err
	}
	s.mode = mode
	return HVACModeChanged{Mode: mode}, nil
}

func (s *HVACSystem) SetRoomPreset(room string, preset ThermostatPreset) (Change, error) {
	th, err := s.findThermostat(room)
	if err != nil {
		return nil, err
	}
	if err = th.setPreset(preset); err != nil {
		return nil, err
	}
	return RoomPresetChanged{Room: room, Preset: preset}, nil
}

func (s *HVACSystem) TurnRoomOn(room string) (Change, error) {
	th, err := s.findThermostat(room)
	if err != nil {
		return nil, err
	}
	th.turnOn()
	return RoomPowerChanged{
		Room: room,
		On:   true,
	}, nil
}

func (s *HVACSystem) TurnRoomOff(room string) (Change, error) {
	th, err := s.findThermostat(room)
	if err != nil {
		return nil, err
	}
	th.turnOff()
	return RoomPowerChanged{
		Room: room,
		On:   false,
	}, nil
}

func (m HVACSystemMode) SupportsSetpoint() bool {
	switch m {
	case HVACSystemModeHeat, HVACSystemModeCool:
		return true
	default:
		return false
	}
}

func (s *HVACSystem) CurrentSetpoint(room string) (float64, error) {
	if err := s.mode.Validate(); err != nil {
		return 0, err
	}

	th, err := s.findThermostat(room)

	if err != nil {
		return 0, err
	}

	if !s.mode.SupportsSetpoint() {
		return 0, ErrCurrentSetPointInvalid
	}

	return th.currentSetpointForMode(s.mode)
}

func (s *HVACSystem) SetTemperature(room string, mode HVACSystemMode, preset ThermostatPreset, temp float64) (Change, error) {
	th, err := s.findThermostat(room)
	if err != nil {
		return nil, err
	}

	if err = th.setTemperature(mode, preset, temp); err != nil {
		return nil, err
	}
	return TemperatureChanged{
		Room:   room,
		Mode:   mode,
		Preset: preset,
		Value:  temp,
	}, nil
}

func (s *HVACSystem) SetCurrentSetPoint(room string, temp float64) (Change, error) {
	th, err := s.findThermostat(room)

	if err != nil {
		return nil, err
	}

	if err = th.setTemperature(s.mode, th.preset, temp); err != nil {
		return nil, err
	}
	return TemperatureChanged{
		Room:   room,
		Mode:   s.mode,
		Preset: th.preset,
		Value:  temp,
	}, nil
}

func (s *HVACSystem) SetHeatComfortTemperature(room string, temp float64) (Change, error) {
	return s.SetTemperature(room, HVACSystemModeHeat, PresetComfort, temp)
}

func (s *HVACSystem) SetHeatEcoTemperature(room string, temp float64) (Change, error) {
	return s.SetTemperature(room, HVACSystemModeHeat, PresetEco, temp)
}

func (s *HVACSystem) SetCoolComfortTemperature(room string, temp float64) (Change, error) {
	return s.SetTemperature(room, HVACSystemModeCool, PresetComfort, temp)
}

func (s *HVACSystem) SetCoolEcoTemperature(room string, temp float64) (Change, error) {
	return s.SetTemperature(room, HVACSystemModeCool, PresetEco, temp)
}

func (s *HVACSystem) findThermostat(room string) (*Thermostat, error) {
	for i := range s.thermostats {
		if strings.EqualFold(s.thermostats[i].room, room) {
			return &s.thermostats[i], nil
		}
	}
	return nil, ErrThermostatNotFound
}
