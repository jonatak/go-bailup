package bootstrap

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jonatak/go-bailup/internal/application"
	"github.com/jonatak/go-bailup/internal/handler/mqtt"
	"github.com/jonatak/go-bailup/internal/infrastructure/bailup"
)

func NewHVACService() (*application.HVACService, error) {
	bailupEmail := os.Getenv("BAILUP_EMAIL")
	bailupPassword := os.Getenv("BAILUP_PASS")
	bailupRegulation := os.Getenv("BAILUP_REGULATION")

	if bailupEmail == "" || bailupPassword == "" || bailupRegulation == "" {
		return nil, InitError
	}

	gateway := bailup.NewGateway(bailupEmail, bailupPassword, bailupRegulation)
	err := gateway.Connect()
	if err != nil {
		return nil, fmt.Errorf("connect HVAC gateway: %w", err)
	}
	return application.NewHVACService(gateway), nil
}

func NewMQTTServer() (*mqtt.Handler, error) {
	host := os.Getenv("MQTT_HOST")
	username := os.Getenv("MQTT_USERNAME")
	password := os.Getenv("MQTT_PASSWORD")
	prefix := os.Getenv("MQTT_TOPIC_PREFIX")
	clientId := os.Getenv("MQTT_CLIENT_ID")
	port, err := strconv.Atoi(os.Getenv("MQTT_PORT"))

	if err != nil {
		return nil, fmt.Errorf("MQTT_PORT invalid: %s", os.Getenv("MQTT_PORT"))
	}

	params := mqtt.ConnectionParams{
		Host:     host,
		Username: username,
		Password: password,
		Port:     port,
		ClientID: clientId,
	}

	handler, err := mqtt.NewMQTTHandler(params, prefix)

	if err != nil {
		return nil, err
	}

	return handler, nil
}
