package domain_test

import (
	"testing"

	"github.com/jonatak/go-bailup/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTemperatureSettingsRejectsInvalidOrderingForMode(t *testing.T) {
	testCases := []struct {
		name    string
		mode    domain.HVACSystemMode
		comfort float64
		eco     float64
	}{
		{
			name:    "heat eco above comfort",
			mode:    domain.HVACSystemModeHeat,
			comfort: 20,
			eco:     21,
		},
		{
			name:    "cool eco below comfort",
			mode:    domain.HVACSystemModeCool,
			comfort: 24,
			eco:     23,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			settings, err := domain.NewTemperatureSettings(tc.mode, tc.comfort, tc.eco)

			require.ErrorIs(t, err, domain.ErrInvalidTemperatureSettingForMode)
			assert.Equal(t, domain.TemperatureSettings{}, settings)
		})
	}
}

func TestNewTemperatureSettingsRejectsTooSmallComfortEcoRange(t *testing.T) {
	testCases := []struct {
		name    string
		mode    domain.HVACSystemMode
		comfort float64
		eco     float64
	}{
		{
			name:    "heat",
			mode:    domain.HVACSystemModeHeat,
			comfort: 20,
			eco:     18.5,
		},
		{
			name:    "cool",
			mode:    domain.HVACSystemModeCool,
			comfort: 24,
			eco:     25.5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			settings, err := domain.NewTemperatureSettings(tc.mode, tc.comfort, tc.eco)

			require.ErrorIs(t, err, domain.ErrInvalidTemperatureRange)
			assert.Equal(t, domain.TemperatureSettings{}, settings)
		})
	}
}

func TestNewTemperatureSettingsRejectsModeWithoutSetpoints(t *testing.T) {
	settings, err := domain.NewTemperatureSettings(domain.HVACSystemModeOff, 20, 18)

	require.ErrorIs(t, err, domain.ErrInvalidTemperatureSettingForMode)
	assert.Equal(t, domain.TemperatureSettings{}, settings)
}
