package mqtt

import (
	"fmt"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jonatak/go-bailup/internal/domain"
)

type subscription struct {
	room       string
	intentChan chan<- intent
	errorChan  chan<- error
}

func (s *subscription) setMode(_ mqtt.Client, msg mqtt.Message) {
	mode := domain.HVACSystemMode(strings.TrimSpace(string(msg.Payload())))

	if err := mode.Validate(); err != nil {
		s.errorChan <- err
		return
	}

	s.intentChan <- setModeIntent{
		mode: mode,
	}
}

func (s *subscription) setTemperature(_ mqtt.Client, msg mqtt.Message) {

	value, err := strconv.ParseFloat(strings.TrimSpace(string(msg.Payload())), 64)
	if err != nil {
		s.errorChan <- err
		return
	}

	s.intentChan <- setTemperatureIntent{
		room:  s.room,
		value: value,
	}
}

func (s *subscription) turnOnOff(_ mqtt.Client, msg mqtt.Message) {
	value := strings.TrimSpace(string(msg.Payload()))
	switch value {
	case "auto":
		s.intentChan <- turnRoomOnIntent{
			room: s.room,
		}
	case "off":
		s.intentChan <- turnRoomOffIntent{
			room: s.room,
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

	s.intentChan <- setPresetIntent{
		room:   s.room,
		preset: preset,
	}
}
