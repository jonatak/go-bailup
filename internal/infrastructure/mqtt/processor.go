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
	service    *application.HVACService
	handler    *Handler
	errorChan  <-chan error
	intentChan <-chan application.Intent
}

func (p *Processor) Run(ctx context.Context) error {
	defer p.handler.Close()
	mqttConnected := false
	for {
		if !mqttConnected {
			if err := p.handler.Connect(); err != nil {
				slog.Error("reconnection failed", "err", err)

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(5 * time.Second):
				}
				continue
			} else {
				mqttConnected = true
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-p.errorChan:
			switch {
			case errors.Is(err, ErrConnectionLost):
				slog.Warn("mqtt connection lost", "err", err)
				mqttConnected = false
			case errors.Is(err, ErrSubscriptionError):
				slog.Error("mqtt subscription failed", "err", err)
				mqttConnected = false
			default:
				slog.Error("mqtt processor error", "err", err)
			}
		case intent := <-p.intentChan:
			ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
			_, err := p.service.ApplyIntent(ctxWithTimeout, intent)
			cancel()
			if err != nil {
				slog.Error("apply intent failed", "intent", fmt.Sprintf("%T", intent), "err", err)
				return err
			}
		}
	}
}
