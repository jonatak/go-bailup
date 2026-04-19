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

func (s *HVACSystem) FindThermostat(room string) (*Thermostat, error) {
	for i := range s.thermostats {
		if strings.EqualFold(s.thermostats[i].room, room) {
			return &s.thermostats[i], nil
		}
	}
	return nil, ErrThermostatNotFound
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

	th, err := s.FindThermostat(room)

	if err != nil {
		return 0, err
	}

	if !s.mode.SupportsSetpoint() {
		return 0, ErrCurrentSetPointInvalid
	}

	return th.currentSetpointForMode(s.mode)
}

func (s *HVACSystem) SetTemperature(room string, mode HVACSystemMode, preset ThermostatPreset, temp float64) error {
	th, err := s.FindThermostat(room)
	if err != nil {
		return err
	}

	return th.setTemperature(mode, preset, temp)
}

func (s *HVACSystem) SetCurrentSetPoint(room string, temp float64) error {
	th, err := s.FindThermostat(room)

	if err != nil {
		return err
	}

	return th.setTemperature(s.mode, th.preset, temp)
}

func (s *HVACSystem) SetHeatComfortTemp(room string, temp float64) error {
	return s.SetTemperature(room, HVACSystemModeHeat, PresetComfort, temp)
}

func (s *HVACSystem) SetHeatEcoTemp(room string, temp float64) error {
	return s.SetTemperature(room, HVACSystemModeHeat, PresetEco, temp)
}

func (s *HVACSystem) SetCoolComfortTemp(room string, temp float64) error {
	return s.SetTemperature(room, HVACSystemModeCool, PresetComfort, temp)
}

func (s *HVACSystem) SetCoolEcoTemp(room string, temp float64) error {
	return s.SetTemperature(room, HVACSystemModeCool, PresetEco, temp)
}
