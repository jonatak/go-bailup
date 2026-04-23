package mqtt

type ConnectionParams struct {
	Host     string
	Username string
	Password string
	ClientID string
	Port     int
}

func (m ConnectionParams) Validate() error {
	if m.Host == "" || m.Username == "" || m.Password == "" || m.ClientID == "" || m.Port <= 0 || m.Port > 65535 {
		return ErrInvalidConnectionParams
	}
	return nil
}
