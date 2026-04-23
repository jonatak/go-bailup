package mqtt

import "github.com/jonatak/go-bailup/internal/application"

type HandlerParams struct {
	Host       string
	Username   string
	Password   string
	ClientID   string
	Port       int
	Prefix     string
	ErrorChan  chan<- error
	IntentChan chan<- application.Intent
}

func (m HandlerParams) Validate() error {
	if m.Host == "" || m.Username == "" || m.Password == "" || m.ClientID == "" || m.Port <= 0 || m.Port > 65535 || m.Prefix == "" || m.ErrorChan == nil || m.IntentChan == nil {
		return ErrInvalidConnectionParams
	}
	return nil
}
