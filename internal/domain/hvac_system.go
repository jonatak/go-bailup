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

func (s *HVACSystem) SetMode(mode HVACSystemMode) error {
	if err := mode.Validate(); err != nil {
		return err
	}
	s.mode = mode
	return nil
}

func (s *HVACSystem) SetRoomPreset(room string, preset ThermostatPreset) error {
	th, err := s.findThermostat(room)
	if s.mode != HVACSystemModeCool && s.mode != HVACSystemModeHeat {
		return ErrInvalidPresetMode
	}
	if err != nil {
		return err
	}
	if err = th.setPreset(preset); err != nil {
		return err
	}
	return nil
}

func (s *HVACSystem) SetRoomPower(room string, on bool) error {
	th, err := s.findThermostat(room)
	if err != nil {
		return err
	}

	if on {
		th.turnOn()
	} else {
		th.turnOff()
	}

	return nil
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
		return 0, ErrCurrentSetpointUnavailable
	}

	return th.currentSetpointForMode(s.mode)
}

func (s *HVACSystem) Setpoint(
	room string,
	mode HVACSystemMode,
	preset ThermostatPreset,
) (float64, error) {
	if err := mode.Validate(); err != nil {
		return 0, err
	}

	if err := preset.Validate(); err != nil {
		return 0, err
	}

	if !mode.SupportsSetpoint() {
		return 0, ErrInvalidTemperatureSettingForMode
	}

	th, err := s.findThermostat(room)
	if err != nil {
		return 0, err
	}

	return th.setpointFor(mode, preset)
}

func (s *HVACSystem) SetTemperature(room string, mode HVACSystemMode, preset ThermostatPreset, temp float64) error {
	th, err := s.findThermostat(room)
	if err != nil {
		return err
	}

	if err = th.setTemperature(mode, preset, temp); err != nil {
		return err
	}
	return nil
}

func (s *HVACSystem) findThermostat(room string) (*Thermostat, error) {
	for i := range s.thermostats {
		if strings.EqualFold(s.thermostats[i].room, room) {
			return &s.thermostats[i], nil
		}
	}
	return nil, ErrThermostatNotFound
}
