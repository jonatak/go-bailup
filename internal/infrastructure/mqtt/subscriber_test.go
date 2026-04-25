package mqtt

import (
	"testing"

	"github.com/jonatak/go-bailup/internal/application"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionTurnOnPublishesIntent(t *testing.T) {
	intents := make(chan application.Intent, 1)
	errors := make(chan error, 1)
	sub := &subscription{
		room:       "living room",
		intentChan: intents,
		errorChan:  errors,
	}

	sub.turnOnOff(nil, testMQTTMessage{payload: []byte("auto")})

	require.Empty(t, errors)
	assert.Equal(t, application.SetRoomPowerIntent{
		Room: "living room",
		On:   true,
	}, <-intents)
}

type testMQTTMessage struct {
	payload []byte
}

func (m testMQTTMessage) Duplicate() bool {
	return false
}

func (m testMQTTMessage) Qos() byte {
	return 0
}

func (m testMQTTMessage) Retained() bool {
	return false
}

func (m testMQTTMessage) Topic() string {
	return ""
}

func (m testMQTTMessage) MessageID() uint16 {
	return 0
}

func (m testMQTTMessage) Payload() []byte {
	return m.payload
}

func (m testMQTTMessage) Ack() {}
