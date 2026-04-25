package mqtt

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jonatak/go-bailup/internal/application"
)

type Processor struct {
	service       *application.HVACService
	handler       *Handler
	mqttConnected bool
}

const refreshInterval = 30 * time.Second

func NewProcessor(handler *Handler, service *application.HVACService) *Processor {
	return &Processor{
		service: service,
		handler: handler,
	}
}

func (p *Processor) Run(ctx context.Context) error {
	defer p.handler.Close()
	timer := time.NewTimer(refreshInterval)
	defer timer.Stop()
	for {

		err := p.ensureMQTTConnected(ctx)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-p.handler.Errors():
			p.handleError(err)
		case intent := <-p.handler.Intents():
			p.handleIntent(ctx, intent)
		case <-timer.C:
			p.refreshState(ctx)
		}
		timer.Reset(refreshInterval)
	}
}
func (p *Processor) handleError(err error) {
	switch {
	case errors.Is(err, ErrConnectionLost):
		slog.Warn("mqtt connection lost", "err", err)
		p.mqttConnected = false
	case errors.Is(err, ErrSubscriptionError):
		slog.Error("mqtt subscription failed", "err", err)
		p.mqttConnected = false
	default:
		slog.Error("mqtt processor error", "err", err)
	}
}

func (p *Processor) ensureMQTTConnected(ctx context.Context) error {
	for !p.mqttConnected {
		if err := p.handler.Connect(); err != nil {
			slog.Error("reconnection failed", "err", err)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(5 * time.Second):
			}
			continue
		}
		p.mqttConnected = true
	}
	return nil
}

func (p *Processor) handleIntent(ctx context.Context, intent application.Intent) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	s, err := p.service.ApplyIntent(ctxWithTimeout, intent)
	cancel()
	if err != nil {
		slog.Error("apply intent failed", "intent", fmt.Sprintf("%T", intent), "err", err)
		return
	}
	fmt.Println(s)
}

func (p *Processor) refreshState(ctx context.Context) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	s, err := p.service.CurrentState(ctxWithTimeout)
	cancel()
	if err != nil {
		slog.Error("state refresh failed", "err", err)
		return
	}
	fmt.Println(s)
}
