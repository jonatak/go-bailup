package mqtt

import (
	"fmt"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jonatak/go-bailup/internal/application"
	"github.com/jonatak/go-bailup/internal/domain"
)

type subscription struct {
	ID               int
	room             string
	thermostatConfig MQTTThermostat
	intentChan       chan<- application.Intent
	errorChan        chan<- error
}

func (s *subscription) setMode(_ mqtt.Client, msg mqtt.Message) {

	mode := ModeToDomain(strings.TrimSpace(string(msg.Payload())))
	if err := mode.Validate(); err != nil {
		s.errorChan <- err
		return
	}

	s.intentChan <- application.SetModeIntent{
		Mode: mode,
	}
}

func (s *subscription) setTemperature(_ mqtt.Client, msg mqtt.Message) {

	value, err := strconv.ParseFloat(strings.TrimSpace(string(msg.Payload())), 64)
	if err != nil {
		s.errorChan <- err
		return
	}

	s.intentChan <- application.SetTemperatureIntent{
		Room:    s.room,
		Preset:  application.TemperaturePresetCurrent,
		Mode:    application.TemperatureModeCurrent,
		Value:   value,
		IsDelta: false,
	}
}

func (s *subscription) turnOnOff(_ mqtt.Client, msg mqtt.Message) {
	value := strings.TrimSpace(string(msg.Payload()))
	switch value {
	case "auto":
		s.intentChan <- application.SetRoomPowerIntent{
			Room: s.room,
			On:   true,
		}
	case "off":
		s.intentChan <- application.SetRoomPowerIntent{
			Room: s.room,
			On:   false,
		}
	default:
		s.errorChan <- fmt.Errorf("received invalid mode %s for room %s", value, s.room)
	}

}

func (s *subscription) setPreset(_ mqtt.Client, msg mqtt.Message) {
	value := strings.TrimSpace(string(msg.Payload()))
	preset := domain.ThermostatPreset(value)
	if err := preset.Validate(); err != nil {
		s.errorChan <- err
		return
	}

	s.intentChan <- application.SetRoomPresetIntent{
		Room:   s.room,
		Preset: preset,
	}
}
