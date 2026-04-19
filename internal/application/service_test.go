package application

import (
	"errors"
	"testing"

	"github.com/jonatak/go-bailup/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHVACServiceCurrentStateReturnsGatewayState(t *testing.T) {
	system := testHVACSystem(t, domain.HVACSystemModeHeat)
	gateway := &fakeHVACSystemGateway{
		state: system,
	}
	service := NewHVACService(gateway)

	got, err := service.CurrentState()

	require.NoError(t, err)
	assert.Same(t, system, got)
	assert.Equal(t, 1, gateway.getStateCalls)
	assert.Equal(t, 0, gateway.applyChangeCalls)
}

func TestHVACServiceSetModeAppliesDomainChange(t *testing.T) {
	updated := testHVACSystem(t, domain.HVACSystemModeCool)
	gateway := &fakeHVACSystemGateway{
		state:        testHVACSystem(t, domain.HVACSystemModeHeat),
		updatedState: updated,
	}
	service := NewHVACService(gateway)

	got, err := service.SetMode(domain.HVACSystemModeCool)

	require.NoError(t, err)
	assert.Same(t, updated, got)
	assert.Equal(t, 1, gateway.getStateCalls)
	assert.Equal(t, 1, gateway.applyChangeCalls)
	assert.Equal(t, domain.HVACModeChanged{
		Mode: domain.HVACSystemModeCool,
	}, gateway.appliedChange)
}

func TestHVACServiceDoesNotApplyChangeWhenStateLoadFails(t *testing.T) {
	wantErr := errors.New("state failed")
	gateway := &fakeHVACSystemGateway{
		getStateErr: wantErr,
	}
	service := NewHVACService(gateway)

	got, err := service.SetMode(domain.HVACSystemModeCool)

	require.ErrorIs(t, err, wantErr)
	assert.Nil(t, got)
	assert.Equal(t, 1, gateway.getStateCalls)
	assert.Equal(t, 0, gateway.applyChangeCalls)
}

func TestHVACServiceDoesNotApplyChangeWhenDomainRejectsCommand(t *testing.T) {
	gateway := &fakeHVACSystemGateway{
		state: testHVACSystem(t, domain.HVACSystemModeHeat),
	}
	service := NewHVACService(gateway)

	got, err := service.SetMode(domain.HVACSystemMode("invalid"))

	require.ErrorIs(t, err, domain.ErrInvalidHVACMode)
	assert.Nil(t, got)
	assert.Equal(t, 1, gateway.getStateCalls)
	assert.Equal(t, 0, gateway.applyChangeCalls)
}

func TestHVACServiceReturnsApplyChangeError(t *testing.T) {
	wantErr := errors.New("apply failed")
	gateway := &fakeHVACSystemGateway{
		state:          testHVACSystem(t, domain.HVACSystemModeHeat),
		applyChangeErr: wantErr,
	}
	service := NewHVACService(gateway)

	got, err := service.SetMode(domain.HVACSystemModeCool)

	require.ErrorIs(t, err, wantErr)
	assert.Nil(t, got)
	assert.Equal(t, 1, gateway.getStateCalls)
	assert.Equal(t, 1, gateway.applyChangeCalls)
}

func TestHVACServiceMethodsApplyExpectedChanges(t *testing.T) {
	testCases := []struct {
		name string
		act  func(*HVACService) (*domain.HVACSystem, error)
		want domain.Change
	}{
		{
			name: "set room preset",
			act: func(service *HVACService) (*domain.HVACSystem, error) {
				return service.SetRoomPreset("Living Room", domain.PresetEco)
			},
			want: domain.RoomPresetChanged{
				Room:   "Living Room",
				Preset: domain.PresetEco,
			},
		},
		{
			name: "turn room on",
			act: func(service *HVACService) (*domain.HVACSystem, error) {
				return service.TurnRoomOn("Living Room")
			},
			want: domain.RoomPowerChanged{
				Room: "Living Room",
				On:   true,
			},
		},
		{
			name: "turn room off",
			act: func(service *HVACService) (*domain.HVACSystem, error) {
				return service.TurnRoomOff("Living Room")
			},
			want: domain.RoomPowerChanged{
				Room: "Living Room",
				On:   false,
			},
		},
		{
			name: "set current setpoint",
			act: func(service *HVACService) (*domain.HVACSystem, error) {
				return service.SetCurrentSetpoint("Living Room", 21)
			},
			want: domain.TemperatureChanged{
				Room:   "Living Room",
				Mode:   domain.HVACSystemModeHeat,
				Preset: domain.PresetComfort,
				Value:  21,
			},
		},
		{
			name: "set temperature",
			act: func(service *HVACService) (*domain.HVACSystem, error) {
				return service.SetTemperature("Living Room", domain.HVACSystemModeCool, domain.PresetEco, 27)
			},
			want: domain.TemperatureChanged{
				Room:   "Living Room",
				Mode:   domain.HVACSystemModeCool,
				Preset: domain.PresetEco,
				Value:  27,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updated := testHVACSystem(t, domain.HVACSystemModeHeat)
			gateway := &fakeHVACSystemGateway{
				state:        testHVACSystem(t, domain.HVACSystemModeHeat),
				updatedState: updated,
			}
			service := NewHVACService(gateway)

			got, err := tc.act(service)

			require.NoError(t, err)
			assert.Same(t, updated, got)
			assert.Equal(t, 1, gateway.getStateCalls)
			assert.Equal(t, 1, gateway.applyChangeCalls)
			assert.Equal(t, tc.want, gateway.appliedChange)
		})
	}
}
