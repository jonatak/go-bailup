package application

import "github.com/jonatak/go-bailup/internal/domain"

type HVACService struct {
	gateway HVACSystemGateway
}

func NewHVACService(gateway HVACSystemGateway) *HVACService {
	return &HVACService{
		gateway: gateway,
	}
}

func (s *HVACService) CurrentState() (*domain.HVACSystem, error) {
	return s.gateway.GetHVACSystemState()
}

func (s *HVACService) SetMode(mode domain.HVACSystemMode) (*domain.HVACSystem, error) {
	system, err := s.gateway.GetHVACSystemState()
	if err != nil {
		return nil, err
	}

	change, err := system.SetMode(mode)
	if err != nil {
		return nil, err
	}

	return s.gateway.ApplyChange(change)
}

func (s *HVACService) SetRoomPreset(
	room string,
	preset domain.ThermostatPreset,
) (*domain.HVACSystem, error) {
	system, err := s.gateway.GetHVACSystemState()
	if err != nil {
		return nil, err
	}

	change, err := system.SetRoomPreset(room, preset)
	if err != nil {
		return nil, err
	}

	return s.gateway.ApplyChange(change)
}

func (s *HVACService) TurnRoomOn(room string) (*domain.HVACSystem, error) {
	system, err := s.gateway.GetHVACSystemState()
	if err != nil {
		return nil, err
	}

	change, err := system.TurnRoomOn(room)
	if err != nil {
		return nil, err
	}

	return s.gateway.ApplyChange(change)
}

func (s *HVACService) TurnRoomOff(room string) (*domain.HVACSystem, error) {
	system, err := s.gateway.GetHVACSystemState()
	if err != nil {
		return nil, err
	}

	change, err := system.TurnRoomOff(room)
	if err != nil {
		return nil, err
	}

	return s.gateway.ApplyChange(change)
}

func (s *HVACService) SetCurrentSetpoint(
	room string,
	temp float64,
) (*domain.HVACSystem, error) {
	system, err := s.gateway.GetHVACSystemState()
	if err != nil {
		return nil, err
	}

	change, err := system.SetCurrentSetpoint(room, temp)
	if err != nil {
		return nil, err
	}

	return s.gateway.ApplyChange(change)
}

func (s *HVACService) SetTemperature(
	room string,
	mode domain.HVACSystemMode,
	preset domain.ThermostatPreset,
	temp float64,
) (*domain.HVACSystem, error) {
	system, err := s.gateway.GetHVACSystemState()
	if err != nil {
		return nil, err
	}

	change, err := system.SetTemperature(room, mode, preset, temp)
	if err != nil {
		return nil, err
	}

	return s.gateway.ApplyChange(change)
}
