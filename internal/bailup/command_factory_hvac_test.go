package bailup_test

import (
	"testing"

	"github.com/jonatak/go-bailup/internal/bailup"
	"github.com/jonatak/go-bailup/internal/bailup/command"
	"github.com/jonatak/go-bailup/internal/bailup/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHVACCommandFactoryValidCommands(t *testing.T) {
	testCase := []struct {
		name string
		mode string
		want model.UCMode
	}{
		{
			name: "UCMode is off",
			mode: "off",
			want: model.UCModeOff,
		},
		{
			name: "UCMode is cool",
			mode: "cool",
			want: model.UCModeCool,
		},
		{
			name: "UCMode is heat",
			mode: "heat",
			want: model.UCModeHeat,
		},
		{
			name: "UCMode is dry",
			mode: "dry",
			want: model.UCModeDry,
		},
		{
			name: "UCMode is fan-only",
			mode: "fan-only",
			want: model.UCModeFanOnly,
		},
	}

	for _, c := range testCase {
		t.Run(c.name, func(t *testing.T) {
			cmd, err := bailup.NewHVACModeCommand(c.mode)

			require.NoError(t, err)

			assert.Equal(t, command.ModeCommand{
				Value: int(c.want),
			}, cmd)
		})
	}
}

func TestHVACCommandFactoryInvalidCommand(t *testing.T) {
	cmd, err := bailup.NewHVACModeCommand("invalid")

	require.Error(t, err)
	assert.Nil(t, cmd)
	assert.Contains(t, err.Error(), "unsupported unit mode")
}
