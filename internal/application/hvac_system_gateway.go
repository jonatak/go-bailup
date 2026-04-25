package application

import (
	"context"

	"github.com/jonatak/go-bailup/internal/domain"
)

type HVACSystemGateway interface {
	Connect(context.Context) error
	GetHVACSystemState(context.Context) (*domain.HVACSystem, error)
	ApplyResolvedIntent(context.Context, ResolvedIntent) (*domain.HVACSystem, error)
}
