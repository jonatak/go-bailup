package application

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/domain"
)

type intentFunc = func(system *domain.HVACSystem) (ResolvedIntent, error)

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

func (s *HVACService) ApplyIntent(intent Intent) (*domain.HVACSystem, error) {
	return s.executeIntent(func(system *domain.HVACSystem) (ResolvedIntent, error) {
		switch i := intent.(type) {
		case SetModeIntent:
			return i, system.SetMode(i.Mode)
		case SetRoomPresetIntent:
			return i, system.SetRoomPreset(i.Room, i.Preset)
		case SetRoomPowerIntent:
			return i, system.SetRoomPower(i.Room, i.On)
		case SetTemperatureIntent:
			resolved, err := resolveTemperatureTarget(system, i)
			if err != nil {
				return nil, err
			}

			return resolved.Intent(), system.SetTemperature(resolved.room, resolved.mode, resolved.preset, resolved.value)
		default:
			return nil, fmt.Errorf("unsupported intent type %T", intent)
		}
	})
}

func (s *HVACService) executeIntent(action intentFunc) (*domain.HVACSystem, error) {
	system, err := s.gateway.GetHVACSystemState()
	if err != nil {
		return nil, err
	}

	intent, err := action(system)
	if err != nil {
		return nil, err
	}

	return s.gateway.ApplyResolvedIntent(intent)
}
