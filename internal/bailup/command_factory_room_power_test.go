package bailup_test

import (
	"testing"

	"github.com/jonatak/go-bailup/internal/bailup"
	"github.com/jonatak/go-bailup/internal/bailup/command"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRoomPowerCommandOn(t *testing.T) {
	state := testState()

	cmd, err := bailup.NewRoomPowerCommand(state, "living room", true)
	require.NoError(t, err)

	assert.Equal(t, command.RoomPowerCommand{
		ThermostatID: 9152,
		On:           true,
	}, cmd)
}

func TestNewRoomPowerCommandOff(t *testing.T) {
	state := testState()

	cmd, err := bailup.NewRoomPowerCommand(state, "living room", false)
	require.NoError(t, err)

	assert.Equal(t, command.RoomPowerCommand{
		ThermostatID: 9152,
		On:           false,
	}, cmd)
}

func TestNewRoomPowerCommandInvalidRoom(t *testing.T) {
	state := testState()

	cmd, err := bailup.NewRoomPowerCommand(state, "living room2", false)
	require.Error(t, err)
	require.Nil(t, cmd)

	assert.Contains(t, err.Error(), "thermostat \"living room2\" not found")
}
