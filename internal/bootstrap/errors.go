package bootstrap

import "errors"

var ErrInit = errors.New("baillconnect email, password, and regulation need to be set")
var ErrMqttInit = errors.New("MQTT host, username, password, topic prefix, client id, and port need to be set")
