package bailup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jonatak/go-bailup/internal/application"
	"github.com/jonatak/go-bailup/internal/domain"
	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/model"
)

var _ application.HVACSystemGateway = (*Gateway)(nil)

const stateCacheTTL = time.Minute

type Gateway struct {
	state         *model.State
	lastRefreshed time.Time
	client        *Bailup
}

func NewGateway(email, password, regulation string) *Gateway {
	return &Gateway{
		client:        NewBailup(email, password, regulation),
		lastRefreshed: time.Now(),
	}
}

func (g *Gateway) Connect(ctx context.Context) error {
	if err := g.client.Connect(ctx); err != nil {
		return fmt.Errorf("%w: %w", application.ErrGatewayUnavailable, err)
	}
	return nil
}

func (g *Gateway) GetHVACSystemState(ctx context.Context) (*domain.HVACSystem, error) {
	err := g.ensureStateLoaded(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStateUnavailable, err)
	}

	system, err := HVACSystemFromState(g.state)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStateUnavailable, err)
	}

	return system, nil
}

func (g *Gateway) ApplyResolvedIntent(ctx context.Context, intent application.ResolvedIntent) (*domain.HVACSystem, error) {

	if err := g.ensureStateLoaded(ctx); err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStateUnavailable, err)
	}

	var result *model.State
	err := g.withReconnect(ctx, func() error {
		cmd, err := CommandFromResolvedIntent(g.state, intent)
		if err != nil {
			return fmt.Errorf("%w: %w", application.ErrChangeRejected, err)
		}
		s, err := g.client.Execute(ctx, cmd)
		if err != nil {
			if errors.Is(err, ErrDisconnected) {
				return fmt.Errorf("%w: %w", application.ErrGatewayUnavailable, err)
			}
			return fmt.Errorf("%w: %w", application.ErrChangeRejected, err)
		}
		result = s
		return nil
	})

	if err != nil {
		return nil, err
	}

	g.setState(result)
	system, err := HVACSystemFromState(result)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", application.ErrStateUnavailable, err)
	}

	return system, nil
}

func (g *Gateway) withReconnect(ctx context.Context, op func() error) error {
	for {
		err := op()
		if err == nil {
			return nil
		}
		if !errors.Is(err, ErrDisconnected) {
			return err
		}

		if connectErr := g.client.Connect(ctx); connectErr != nil {
			select {
			case <-time.After(1 * time.Second):
			case <-ctx.Done():
				return application.ErrGatewayUnavailable
			}
			continue
		}
	}
}

func (g *Gateway) ensureStateLoaded(ctx context.Context) error {
	if !g.stateExpired(time.Now()) {
		return nil
	}

	return g.withReconnect(ctx, func() error {
		s, err := g.client.GetState(ctx)
		if err != nil {
			return err
		}
		g.setState(s)
		return nil
	})
}

func (g *Gateway) stateExpired(now time.Time) bool {
	return g.state == nil || now.After(g.lastRefreshed.Add(stateCacheTTL))
}

func (g *Gateway) setState(s *model.State) {
	g.state = s
	g.lastRefreshed = time.Now()
}
