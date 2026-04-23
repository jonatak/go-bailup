package mqtt

import (
	"fmt"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jonatak/go-bailup/internal/application"
	"github.com/jonatak/go-bailup/internal/domain"
)

type Handler struct {
	client        mqtt.Client
	prefix        string
	subscriptions []*subscription
	errorChan     chan<- error
}

func NewMQTTHandler(param HandlerParams, system *domain.HVACSystem) (*Handler, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	mqttContext := &Handler{
		client:    nil,
		prefix:    param.Prefix,
		errorChan: param.ErrorChan,
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", param.Host, param.Port))
	opts.SetClientID(param.ClientID)
	opts.SetUsername(param.Username)
	opts.SetPassword(param.Password)
	opts.SetDefaultPublishHandler(mqttContext.messageHandler)
	opts.OnConnect = mqttContext.connectionHandler
	opts.OnConnectionLost = mqttContext.connectionLostHandler

	mqttContext.client = mqtt.NewClient(opts)

	mqttContext.registerSubscription(system, param.IntentChan)

	return mqttContext, nil
}

func (m *Handler) Connect() error {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *Handler) Close() {
	m.client.Disconnect(100)
}

func (m *Handler) registerSubscription(system *domain.HVACSystem, intentChan chan<- application.Intent) {
	th := system.Thermostats()
	subscriber := make([]*subscription, 0, len(th)+1)
	subscriber = append(subscriber, &subscription{
		room:       "general",
		intentChan: intentChan,
		errorChan:  m.errorChan,
	})
	for _, t := range th {
		roomName := strings.ToLower(t.Room())
		subscriber = append(subscriber, &subscription{
			room:       roomName,
			intentChan: intentChan,
			errorChan:  m.errorChan,
		})
	}
	m.subscriptions = subscriber
}

func (m *Handler) subscribe() error {

	for _, s := range m.subscriptions {

		switch s.room {
		case "general":
			if token := m.client.Subscribe(fmt.Sprintf("%s/general/mode/set", m.prefix), byte(0), s.setMode); token.Wait() && token.Error() != nil {
				return token.Error()
			}
		default:
			if token := m.client.Subscribe(fmt.Sprintf("%s/%s/preset_mode/set", m.prefix, s.room), byte(0), s.setPreset); token.Wait() && token.Error() != nil {
				return token.Error()
			}
			if token := m.client.Subscribe(fmt.Sprintf("%s/%s/mode/set", m.prefix, s.room), byte(0), s.turnOnOff); token.Wait() && token.Error() != nil {
				return token.Error()
			}
			if token := m.client.Subscribe(fmt.Sprintf("%s/%s/temperature/set", m.prefix, s.room), byte(0), s.setTemperature); token.Wait() && token.Error() != nil {
				return token.Error()
			}
		}
	}

	return nil
}

func (m *Handler) messageHandler(_ mqtt.Client, msg mqtt.Message) {
	log.Printf("received message")
}

func (m *Handler) connectionHandler(client mqtt.Client) {
	if err := m.subscribe(); err != nil && m.errorChan != nil {
		m.errorChan <- fmt.Errorf("%w: %w", ErrSubscriptionError, err)
	}
}

func (m *Handler) connectionLostHandler(client mqtt.Client, err error) {
	if m.errorChan != nil {
		m.errorChan <- fmt.Errorf("%w: %w", ErrConnectionLost, err)
	}
}
