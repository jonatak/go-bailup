package mqtt

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jonatak/go-bailup/internal/application"
	"github.com/jonatak/go-bailup/internal/domain"
)

type job interface {
	isJob()
}

type intentJob struct {
	intent application.Intent
}

type refreshJob struct{}

func (intentJob) isJob()  {}
func (refreshJob) isJob() {}

type result struct {
	state *domain.HVACSystem
	err   error
}

type Processor struct {
	service       *application.HVACService
	handler       *Handler
	mqttConnected bool
}

const refreshInterval = 60 * time.Second

func NewProcessor(handler *Handler, service *application.HVACService) *Processor {
	return &Processor{
		service: service,
		handler: handler,
	}
}

func (p *Processor) Run(ctx context.Context) error {
	defer p.handler.Close()
	timer := time.NewTimer(refreshInterval)
	jobCh := make(chan job, 10)
	resultCh := make(chan result)
	defer close(jobCh)

	go p.StartWorker(ctx, jobCh, resultCh)
	defer timer.Stop()

	slog.Info("Processor started")
	jobCh <- refreshJob{}
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
		case res := <-resultCh:
			if res.err != nil {
				slog.Error(res.err.Error())
				continue
			}
			if err := p.handler.PublishState(res.state); err != nil {
				p.handleError(err)
				continue
			}
		case intent := <-p.handler.Intents():
			slog.Info("received intent", "intent", intent)
			if len(jobCh) == cap(jobCh) {
				slog.Info("mqtt command dropped, worker queue is full")
				continue
			}
			jobCh <- intentJob{
				intent: intent,
			}
		case <-timer.C:
			if len(jobCh) == 0 {
				jobCh <- refreshJob{}
			}
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
	case errors.Is(err, ErrRegistryError):
		slog.Error("mqtt registry failed", "err", err)
		p.mqttConnected = false
	case errors.Is(err, ErrPublishError):
		slog.Error("state publishing failed", "err", err)
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

func (p *Processor) handleIntent(ctx context.Context, intent application.Intent) result {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	s, err := p.service.ApplyIntent(ctxWithTimeout, intent)
	cancel()
	if err != nil {
		return result{
			err: fmt.Errorf("apply intent failed, intent: %T, error: %w", intent, err),
		}
	}
	return result{
		state: s,
	}
}

func (p *Processor) refreshState(ctx context.Context) result {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	s, err := p.service.CurrentState(ctxWithTimeout)
	cancel()
	if err != nil {
		return result{
			err: fmt.Errorf("refresh state failed error: %w", err),
		}
	}
	return result{
		state: s,
	}
}

func (p *Processor) StartWorker(ctx context.Context, jobs <-chan job, results chan<- result) {
	for j := range jobs {
		switch j := j.(type) {
		case intentJob:
			results <- p.handleIntent(ctx, j.intent)
		case refreshJob:
			results <- p.refreshState(ctx)
		}
	}
}
