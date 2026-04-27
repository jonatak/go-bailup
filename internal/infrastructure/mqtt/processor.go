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
			slog.Info("processor case: context done", "err", ctx.Err())
			return ctx.Err()
		case err := <-p.handler.Errors():
			slog.Info("processor case: mqtt error", "err", err)
			p.handleError(err)
		case res := <-resultCh:
			if p.handleWorkerResult(res) {
				timer.Reset(refreshInterval)
			}
		case intent := <-p.handler.Intents():
			if !p.handleIntentMsg(intent, jobCh) {
				timer.Reset(refreshInterval)
			}
		case <-timer.C:
			p.handleInactivityTimer(jobCh)
			timer.Reset(refreshInterval)
		}
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

func (p *Processor) handleIntentWorker(ctx context.Context, intent application.Intent) result {
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
			results <- p.handleIntentWorker(ctx, j.intent)
		case refreshJob:
			results <- p.refreshState(ctx)
		}
	}
}

func (p *Processor) handleWorkerResult(res result) bool {
	slog.Info("processor case: worker result", "has_state", res.state != nil, "err", res.err)

	if res.err != nil {
		slog.Error(res.err.Error())
		return false
	}

	if err := p.handler.PublishState(res.state); err != nil {
		p.handleError(err)
		return false
	}

	return true
}

func (p *Processor) handleInactivityTimer(jobCh chan<- job) {
	slog.Info("processor case: inactivity timer")
	if len(jobCh) != 0 {
		slog.Info("refresh skipped, worker queue is not empty", "queue_len", len(jobCh), "queue_cap", cap(jobCh))
		return
	}
	jobCh <- refreshJob{}
}

func (p *Processor) handleIntentMsg(intent application.Intent, jobCh chan<- job) bool {
	slog.Info("processor case: intent", "intent", intent, "queue_len", len(jobCh), "queue_cap", cap(jobCh))
	if len(jobCh) == cap(jobCh) {
		slog.Info("mqtt command dropped, worker queue is full")
		return false
	}
	jobCh <- intentJob{
		intent: intent,
	}
	return true
}
