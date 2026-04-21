package mqtt

import "errors"

var (
	ErrInvalidConnectionParams = errors.New("invalid MQTT connection params")
	ErrInvalidTopic            = errors.New("invalid MQTT topic")
)
