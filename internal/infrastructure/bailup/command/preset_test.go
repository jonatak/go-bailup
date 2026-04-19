package command_test

import (
	"testing"

	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/command"
	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPresetCommandComfort(t *testing.T) {
	state := testState()

	cmd, err := command.NewPresetCommand(state, "living room", "comfort")
	require.NoError(t, err)

	assert.Equal(t, command.Preset{
		ThermostatID: 9152,
		Value:        int(model.ThModeComfort),
	}, cmd)
}

func TestNewPresetCommandEco(t *testing.T) {
	state := testState()

	cmd, err := command.NewPresetCommand(state, "bedroom", "eco")
	require.NoError(t, err)

	assert.Equal(t, command.Preset{
		ThermostatID: 9154,
		Value:        int(model.ThModeEco),
	}, cmd)
}

func TestNewPresetCommandInvalidRoom(t *testing.T) {
	state := testState()

	cmd, err := command.NewPresetCommand(state, "kitchen", "comfort")

	require.Error(t, err)
	assert.Nil(t, cmd)
	assert.Contains(t, err.Error(), "thermostat \"kitchen\" not found")
}

func TestNewPresetCommandInvalidPreset(t *testing.T) {
	state := testState()

	cmd, err := command.NewPresetCommand(state, "living room", "away")

	require.Error(t, err)
	assert.Nil(t, cmd)
	assert.Contains(t, err.Error(), "unsupported thermostat mode")
}
