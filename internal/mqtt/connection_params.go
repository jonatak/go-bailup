package mqtt

type HandlerParams struct {
	Host     string
	Username string
	Password string
	ClientID string
	Port     int
	Prefix   string
}

func (m HandlerParams) Validate() error {
	if m.Host == "" || m.Username == "" || m.Password == "" || m.ClientID == "" || m.Port <= 0 || m.Port > 65535 || m.Prefix == "" {
		return ErrInvalidConnectionParams
	}
	return nil
}
