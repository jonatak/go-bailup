package application

import (
	"context"
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

func (s *HVACService) CurrentState(ctx context.Context) (*domain.HVACSystem, error) {
	return s.gateway.GetHVACSystemState(ctx)
}

func (s *HVACService) ApplyIntent(ctx context.Context, intent Intent) (*domain.HVACSystem, error) {
	return s.executeIntent(ctx, func(system *domain.HVACSystem) (ResolvedIntent, error) {
		switch i := intent.(type) {
		case SetModeIntent:
			return i, system.SetMode(i.Mode)
		case SetRoomPresetIntent:
			return i, system.SetRoomPreset(i.Room, i.Preset)
		case SetRoomPowerIntent:
			if i.On && system.Mode() == domain.HVACSystemModeOff {
				return nil, nil
			}
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

func (s *HVACService) executeIntent(ctx context.Context, action intentFunc) (*domain.HVACSystem, error) {
	system, err := s.gateway.GetHVACSystemState(ctx)
	if err != nil {
		return nil, err
	}

	intent, err := action(system)
	if err != nil {
		return nil, err
	}
	if intent == nil {
		return system, nil
	}

	return s.gateway.ApplyResolvedIntent(ctx, intent)
}
