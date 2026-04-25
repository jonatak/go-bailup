package bootstrap

import "errors"

var InitError = errors.New("env var BAILUP_EMAIL, BAILUP_PASS, BAILUP_REGULATION need to be set")
var MqttInitError = errors.New("env var MQTT_HOST, MQTT_USERNAME, MQTT_PASSWORD, MQTT_TOPIC_PREFIX, MQTT_CLIENT_ID, MQTT_PORT need to be set")
