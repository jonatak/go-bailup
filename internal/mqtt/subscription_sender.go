package mqtt

import (
	"log/slog"

	"github.com/jonatak/baillconnect-to-mqtt/internal/application"
)

type subscriptionSender struct {
	intentChan chan<- application.Intent
	errorChan  chan<- error
}

func (s subscriptionSender) sendIntent(i application.Intent) {
	select {
	case s.intentChan <- i:
	default:
		slog.Info("subscriber couldn't send intent, intentChan is full.")
	}
}

func (s subscriptionSender) sendError(err error) {
	select {
	case s.errorChan <- err:
	default:
		slog.Error("subscriber couldn't send error, errorChan is full.", "error", err)
	}
}
