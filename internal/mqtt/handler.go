package mqtt

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jonatak/baillconnect-to-mqtt/internal/application"
	"github.com/jonatak/baillconnect-to-mqtt/internal/domain"
)

const discoveryTopic = "homeassistant"

type Handler struct {
	client      mqtt.Client
	prefix      string
	general     *generalSubscription
	thermostats map[int]*subscription
	errorChan   chan error
	intentChan  chan application.Intent
}

func NewMQTTHandler(param HandlerParams, system *domain.HVACSystem) (*Handler, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	mqttContext := &Handler{
		client:     nil,
		prefix:     param.Prefix,
		errorChan:  make(chan error, 10),
		intentChan: make(chan application.Intent, 10),
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

	mqttContext.registerSubscription(system, mqttContext.intentChan)

	return mqttContext, nil
}

func (m *Handler) Errors() <-chan error {
	return m.errorChan
}

func (m *Handler) Intents() <-chan application.Intent {
	return m.intentChan
}

func (m *Handler) Connect() error {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *Handler) Close() {
	if m.client.IsConnected() {
		m.client.Disconnect(100)
	}
}

func (m *Handler) PublishState(system *domain.HVACSystem) error {

	tempTotal := 0.0
	values := make(map[string]string)

	systemMode := system.Mode()

	for _, t := range system.Thermostats() {
		tempTotal += t.Temperature()
		s := m.thermostats[t.ID()]
		c := s.thermostatConfig
		b := s.batteryConfig
		mode := "off"

		values[c.PresetModeStateTopic] = ""
		values[c.TemperatureStateTopic] = ""
		values[b.StateTopic] = "100"
		values[c.CurrentTemperatureTopic] = formatFloat(t.Temperature())

		if systemMode != domain.HVACSystemModeOff && t.IsOn() {
			mode = "auto"
		}

		values[c.ModeStateTopic] = mode
		setPoint, err := system.CurrentSetpoint(t.Room())
		if err == nil {
			values[c.TemperatureStateTopic] = formatFloat(setPoint)
			values[c.PresetModeStateTopic] = PresetFromDomain(t.Preset())
		}

		if t.IsBatteryLow() {
			values[b.StateTopic] = "1"
		}

		action, err := t.Action(systemMode)
		if err != nil {
			return err
		}
		if systemMode == domain.HVACSystemModeOff {
			action = domain.ThermostatActionOff
		}

		values[c.ActionTopic] = string(action)
	}

	values[m.general.thermostatConfig.ModeStateTopic] = ModeFromDomain(systemMode)
	values[m.general.thermostatConfig.CurrentTemperatureTopic] = formatFloat(tempTotal / float64(len(system.Thermostats())))

	for t, v := range values {
		if err := m.publishHelperState(t, v); err != nil {
			return err
		}
	}

	return nil
}

func (m *Handler) publishHelperState(topic string, value string) error {
	if token := m.client.Publish(topic, byte(0), false, value); token.Wait() && token.Error() != nil {
		return fmt.Errorf("%w: %v", ErrPublishError, token.Error())
	}
	return nil
}

func (m *Handler) registerSubscription(system *domain.HVACSystem, intentChan chan<- application.Intent) {
	th := system.Thermostats()
	subscriber := make(map[int]*subscription)
	m.thermostats = subscriber
	sender := subscriptionSender{
		intentChan: intentChan,
		errorChan:  m.errorChan,
	}
	m.general = &generalSubscription{
		subscriptionSender: sender,
		room:               "general",
		thermostatConfig:   ThermostatGeneralFromDomain(m.prefix),
	}
	for _, t := range th {
		roomName := strings.ToLower(t.Room())
		m.thermostats[t.ID()] = &subscription{
			subscriptionSender: sender,
			ID:                 t.ID(),
			room:               roomName,
			thermostatConfig:   ThermostatFromDomain(t, m.prefix),
			batteryConfig:      BatteryFromThermostatDomain(t, m.prefix),
		}
	}
}

func (m *Handler) registerDiscovery() error {

	j, err := json.Marshal(m.general.thermostatConfig)
	if err != nil {
		return err
	}

	if token := m.client.Publish(fmt.Sprintf("%s/climate/general/config", discoveryTopic), byte(0), true, j); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for _, t := range m.thermostats {
		j, err := json.Marshal(t.thermostatConfig)
		if err != nil {
			return err
		}

		b, err := json.Marshal(t.batteryConfig)
		if err != nil {
			return err
		}

		if token := m.client.Publish(fmt.Sprintf("%s/climate/th_%d/config", discoveryTopic, t.ID), byte(0), true, j); token.Wait() && token.Error() != nil {
			return token.Error()
		}
		if token := m.client.Publish(fmt.Sprintf("%s/sensor/th_%d_battery/config", discoveryTopic, t.ID), byte(0), true, b); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}
	return nil
}

func (m *Handler) subscribe() error {
	if token := m.client.Subscribe(m.general.thermostatConfig.ModeCommandTopic, byte(0), m.general.setMode); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for _, s := range m.thermostats {
		if token := m.client.Subscribe(s.thermostatConfig.PresetModeCommandTopic, byte(0), s.setPreset); token.Wait() && token.Error() != nil {
			return token.Error()
		}
		if token := m.client.Subscribe(s.thermostatConfig.ModeCommandTopic, byte(0), s.turnOnOff); token.Wait() && token.Error() != nil {
			return token.Error()
		}
		if token := m.client.Subscribe(s.thermostatConfig.TemperatureCommandTopic, byte(0), s.setTemperature); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	return nil
}

func (m *Handler) messageHandler(_ mqtt.Client, msg mqtt.Message) {
	m.sendError(fmt.Errorf("received unhandled message from topic: %s", msg.Topic()))
}

func (m *Handler) connectionHandler(client mqtt.Client) {
	if err := m.subscribe(); err != nil {
		m.sendError(fmt.Errorf("%w: %w", ErrSubscriptionError, err))
		return
	}
	if err := m.registerDiscovery(); err != nil {
		m.sendError(fmt.Errorf("%w: %w", ErrRegistryError, err))
	}
}

func (m *Handler) connectionLostHandler(client mqtt.Client, err error) {
	if m.errorChan != nil {
		m.sendError(fmt.Errorf("%w: %w", ErrConnectionLost, err))
	}
}

func (m *Handler) sendError(err error) {
	select {
	case m.errorChan <- err:
	default:
		slog.Error("handler couldn't send error, errorChan is full.", "error", err)
	}
}
