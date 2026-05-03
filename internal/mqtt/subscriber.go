package mqtt

import (
	"fmt"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jonatak/baillconnect-to-mqtt/internal/application"
	"github.com/jonatak/baillconnect-to-mqtt/internal/domain"
)

type subscription struct {
	subscriptionSender

	ID               int
	room             string
	thermostatConfig MQTTThermostat
	batteryConfig    MQTTBatterySensor
}

func (s *subscription) setTemperature(_ mqtt.Client, msg mqtt.Message) {

	value, err := strconv.ParseFloat(strings.TrimSpace(string(msg.Payload())), 64)
	if err != nil {
		s.sendError(err)
		return
	}

	s.sendIntent(application.SetTemperatureIntent{
		Room:    s.room,
		Preset:  application.TemperaturePresetCurrent,
		Mode:    application.TemperatureModeCurrent,
		Value:   value,
		IsDelta: false,
	})
}

func (s *subscription) turnOnOff(_ mqtt.Client, msg mqtt.Message) {
	value := strings.TrimSpace(string(msg.Payload()))
	switch value {
	case "auto":
		s.sendIntent(application.SetRoomPowerIntent{
			Room: s.room,
			On:   true,
		})
	case "off":
		s.sendIntent(application.SetRoomPowerIntent{
			Room: s.room,
			On:   false,
		})
	default:
		s.sendError(fmt.Errorf("received invalid mode %s for room %s", value, s.room))
	}

}

func (s *subscription) setPreset(_ mqtt.Client, msg mqtt.Message) {
	value := strings.TrimSpace(string(msg.Payload()))
	preset := domain.ThermostatPreset(value)
	if err := preset.Validate(); err != nil {
		s.sendError(err)
		return
	}

	s.sendIntent(application.SetRoomPresetIntent{
		Room:   s.room,
		Preset: preset,
	})
}
