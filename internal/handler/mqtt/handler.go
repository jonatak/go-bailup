package mqtt

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Handler struct {
	client mqtt.Client
	prefix string
}

func NewMQTTHandler(param ConnectionParams, prefix string) (*Handler, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	mqttContext := &Handler{
		client: nil,
		prefix: prefix,
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
	return mqttContext, nil
}

func (m *Handler) messageHandler(_ mqtt.Client, msg mqtt.Message) {
	log.Printf("received message")
}

func (m *Handler) connectionHandler(client mqtt.Client) {
	log.Println("Connected")
}

func (m *Handler) connectionLostHandler(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}
