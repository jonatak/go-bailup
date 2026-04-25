package application

import (
	"context"
	"testing"

	"github.com/jonatak/go-bailup/internal/domain"
	"github.com/stretchr/testify/require"
)

type fakeHVACSystemGateway struct {
	state          *domain.HVACSystem
	updatedState   *domain.HVACSystem
	getStateErr    error
	applyIntentErr error

	getStateCalls    int
	applyIntentCalls int
	appliedIntent    ResolvedIntent
}

func (f *fakeHVACSystemGateway) Connect(context.Context) error {
	return nil
}

func (f *fakeHVACSystemGateway) GetHVACSystemState(context.Context) (*domain.HVACSystem, error) {
	f.getStateCalls++
	if f.getStateErr != nil {
		return nil, f.getStateErr
	}
	return f.state, nil
}

func (f *fakeHVACSystemGateway) ApplyResolvedIntent(_ context.Context, intent ResolvedIntent) (*domain.HVACSystem, error) {
	f.applyIntentCalls++
	f.appliedIntent = intent
	if f.applyIntentErr != nil {
		return nil, f.applyIntentErr
	}
	return f.updatedState, nil
}

func testHVACSystem(t *testing.T, mode domain.HVACSystemMode) *domain.HVACSystem {
	t.Helper()

	heatSetting, err := domain.NewTemperatureSettings(domain.HVACSystemModeHeat, 20, 18)
	require.NoError(t, err)
	coolSetting, err := domain.NewTemperatureSettings(domain.HVACSystemModeCool, 24, 26)
	require.NoError(t, err)
	thermostat, err := domain.NewThermostat(
		1,
		"Living Room",
		domain.PresetComfort,
		false,
		false,
		heatSetting,
		coolSetting,
	)
	require.NoError(t, err)

	system, err := domain.NewHVACSystem(mode, []domain.Thermostat{thermostat})
	require.NoError(t, err)

	return system
}
