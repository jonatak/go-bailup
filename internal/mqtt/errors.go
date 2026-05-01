package mqtt

import "errors"

var (
	ErrInvalidConnectionParams = errors.New("invalid MQTT connection params")
	ErrInvalidTopic            = errors.New("invalid MQTT topic")
	ErrConnectionLost          = errors.New("mqtt connection lost")
	ErrSubscriptionError       = errors.New("subscription error")
	ErrRegistryError           = errors.New("registration error")
	ErrPublishError            = errors.New("couldn't publish message")
)
