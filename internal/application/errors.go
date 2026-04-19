package application

import "errors"

var (
	ErrGatewayUnavailable = errors.New("gateway unavailable")
	ErrStateUnavailable   = errors.New("hvac system state unavailable")
	ErrChangeRejected     = errors.New("hvac system change rejected")
)
