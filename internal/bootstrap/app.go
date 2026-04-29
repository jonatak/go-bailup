package bootstrap

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jonatak/go-bailup/internal/application"
	"github.com/jonatak/go-bailup/internal/infrastructure/bailup"
	"github.com/jonatak/go-bailup/internal/infrastructure/mqtt"
)

func NewHVACService() (*application.HVACService, error) {
	bailupEmail := os.Getenv("BAILUP_EMAIL")
	bailupPassword := os.Getenv("BAILUP_PASS")
	bailupRegulation := os.Getenv("BAILUP_REGULATION")

	if bailupEmail == "" || bailupPassword == "" || bailupRegulation == "" {
		return nil, ErrInit
	}

	gateway := bailup.NewGateway(bailupEmail, bailupPassword, bailupRegulation)
	err := gateway.Connect(context.Background())
	if err != nil {
		return nil, fmt.Errorf("connect HVAC gateway: %w", err)
	}
	return application.NewHVACService(gateway), nil
}

func NewMQTTServer(
	system *application.HVACService,
) (*mqtt.Processor, error) {

	state, err := system.CurrentState(context.Background())
	if err != nil {
		return nil, err
	}

	host := os.Getenv("MQTT_HOST")
	username := os.Getenv("MQTT_USERNAME")
	password := os.Getenv("MQTT_PASSWORD")
	prefix := os.Getenv("MQTT_TOPIC_PREFIX")
	clientId := os.Getenv("MQTT_CLIENT_ID")
	port, err := strconv.Atoi(os.Getenv("MQTT_PORT"))

	if err != nil {
		return nil, fmt.Errorf("MQTT_PORT invalid: %s", os.Getenv("MQTT_PORT"))
	}

	if host == "" || username == "" || password == "" || prefix == "" || clientId == "" {
		return nil, ErrMqttInit
	}

	params := mqtt.HandlerParams{
		Host:     host,
		Username: username,
		Password: password,
		Port:     port,
		ClientID: clientId,
		Prefix:   prefix,
	}

	handler, err := mqtt.NewMQTTHandler(params, state)

	if err != nil {
		return nil, err
	}

	return mqtt.NewProcessor(handler, system), nil
}
