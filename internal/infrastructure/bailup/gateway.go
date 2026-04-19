package bailup

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/application"
	"github.com/jonatak/go-bailup/internal/domain"
	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/model"
)

var _ application.HVACSystemGateway = (*Gateway)(nil)

type Gateway struct {
	state  *model.State
	client *Bailup
}

func NewGateway(email, password, regulation string) *Gateway {
	return &Gateway{
		client: NewBailup(email, password, regulation),
	}
}

func (g *Gateway) Connect() error {
	if err := g.client.Connect(); err != nil {
		return fmt.Errorf("%w: %w", application.ErrGatewayUnavailable, err)
	}
	return nil
}

func (g *Gateway) GetHVACSystemState() (*domain.HVACSystem, error) {
	err := g.ensureStateLoaded()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStateUnavailable, err)
	}

	system, err := HVACSystemFromState(g.state)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStateUnavailable, err)
	}

	return system, nil
}

func (g *Gateway) ApplyChange(change domain.Change) (*domain.HVACSystem, error) {
	err := g.ensureStateLoaded()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStateUnavailable, err)
	}

	cmd, err := CommandFromChange(g.state, change)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrChangeRejected, err)
	}

	s, err := g.client.Execute(cmd)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrChangeRejected, err)
	}

	g.state = s
	system, err := HVACSystemFromState(s)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStateUnavailable, err)
	}

	return system, nil
}

func (g *Gateway) ensureStateLoaded() error {
	if g.state == nil {
		s, err := g.client.GetState()
		if err != nil {
			return err
		}
		g.state = s
	}
	return nil
}
