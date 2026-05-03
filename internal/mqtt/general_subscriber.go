package mqtt

import (
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jonatak/baillconnect-to-mqtt/internal/application"
)

type generalSubscription struct {
	subscriptionSender

	room             string
	thermostatConfig MQTTGeneralThermostat
}

func (s *generalSubscription) setMode(_ mqtt.Client, msg mqtt.Message) {

	mode := ModeToDomain(strings.TrimSpace(string(msg.Payload())))
	if err := mode.Validate(); err != nil {
		s.sendError(err)
		return
	}

	s.sendIntent(application.SetModeIntent{
		Mode: mode,
	})
}
