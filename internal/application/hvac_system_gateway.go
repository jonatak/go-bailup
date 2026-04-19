package application

import "github.com/jonatak/go-bailup/internal/domain"

type HVACSystemGateway interface {
	Connect() error
	GetHVACSystemState() (*domain.HVACSystem, error)
	ApplyChange(domain.Change) (*domain.HVACSystem, error)
}
