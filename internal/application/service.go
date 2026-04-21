package application

import "github.com/jonatak/go-bailup/internal/domain"

type changeFunc = func(system *domain.HVACSystem) (domain.Change, error)

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
	return s.executeChange(func(system *domain.HVACSystem) (domain.Change, error) {
		return system.SetMode(mode)
	})
}

func (s *HVACService) SetRoomPreset(
	room string,
	preset domain.ThermostatPreset,
) (*domain.HVACSystem, error) {
	return s.executeChange(func(system *domain.HVACSystem) (domain.Change, error) {
		return system.SetRoomPreset(room, preset)
	})
}

func (s *HVACService) TurnRoomOn(room string) (*domain.HVACSystem, error) {
	return s.executeChange(func(system *domain.HVACSystem) (domain.Change, error) {
		return system.TurnRoomOn(room)
	})
}

func (s *HVACService) TurnRoomOff(room string) (*domain.HVACSystem, error) {
	return s.executeChange(func(system *domain.HVACSystem) (domain.Change, error) {
		return system.TurnRoomOff(room)
	})
}

func (s *HVACService) SetCurrentSetpoint(
	room string,
	temp float64,
) (*domain.HVACSystem, error) {
	return s.executeChange(func(system *domain.HVACSystem) (domain.Change, error) {
		return system.SetCurrentSetpoint(room, temp)
	})
}

func (s *HVACService) SetTemperature(
	room string,
	mode domain.HVACSystemMode,
	preset domain.ThermostatPreset,
	temp float64,
) (*domain.HVACSystem, error) {
	return s.executeChange(func(system *domain.HVACSystem) (domain.Change, error) {
		return system.SetTemperature(room, mode, preset, temp)
	})
}

func (s *HVACService) executeChange(action changeFunc) (*domain.HVACSystem, error) {
	system, err := s.gateway.GetHVACSystemState()
	if err != nil {
		return nil, err
	}

	change, err := action(system)
	if err != nil {
		return nil, err
	}

	return s.gateway.ApplyChange(change)
}
