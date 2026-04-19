package application

import (
	"testing"

	"github.com/jonatak/go-bailup/internal/domain"
	"github.com/stretchr/testify/require"
)

type fakeHVACSystemGateway struct {
	state          *domain.HVACSystem
	updatedState   *domain.HVACSystem
	getStateErr    error
	applyChangeErr error

	getStateCalls    int
	applyChangeCalls int
	appliedChange    domain.Change
}

func (f *fakeHVACSystemGateway) Connect() error {
	return nil
}

func (f *fakeHVACSystemGateway) GetHVACSystemState() (*domain.HVACSystem, error) {
	f.getStateCalls++
	if f.getStateErr != nil {
		return nil, f.getStateErr
	}
	return f.state, nil
}

func (f *fakeHVACSystemGateway) ApplyChange(change domain.Change) (*domain.HVACSystem, error) {
	f.applyChangeCalls++
	f.appliedChange = change
	if f.applyChangeErr != nil {
		return nil, f.applyChangeErr
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
