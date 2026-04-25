package application

import (
	"context"
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

	got, err := service.CurrentState(context.Background())

	require.NoError(t, err)
	assert.Same(t, system, got)
	assert.Equal(t, 1, gateway.getStateCalls)
	assert.Equal(t, 0, gateway.applyIntentCalls)
}

func TestHVACServiceApplyIntentSetMode(t *testing.T) {
	updated := testHVACSystem(t, domain.HVACSystemModeCool)
	gateway := &fakeHVACSystemGateway{
		state:        testHVACSystem(t, domain.HVACSystemModeHeat),
		updatedState: updated,
	}
	service := NewHVACService(gateway)

	got, err := service.ApplyIntent(context.Background(), SetModeIntent{Mode: domain.HVACSystemModeCool})

	require.NoError(t, err)
	assert.Same(t, updated, got)
	assert.Equal(t, 1, gateway.getStateCalls)
	assert.Equal(t, 1, gateway.applyIntentCalls)
	assert.Equal(t, SetModeIntent{
		Mode: domain.HVACSystemModeCool,
	}, gateway.appliedIntent)
}

func TestHVACServiceApplyIntentDoesNotApplyWhenStateLoadFails(t *testing.T) {
	wantErr := errors.New("state failed")
	gateway := &fakeHVACSystemGateway{
		getStateErr: wantErr,
	}
	service := NewHVACService(gateway)

	got, err := service.ApplyIntent(context.Background(), SetModeIntent{Mode: domain.HVACSystemModeCool})

	require.ErrorIs(t, err, wantErr)
	assert.Nil(t, got)
	assert.Equal(t, 1, gateway.getStateCalls)
	assert.Equal(t, 0, gateway.applyIntentCalls)
}

func TestHVACServiceApplyIntentDoesNotApplyWhenDomainRejectsCommand(t *testing.T) {
	gateway := &fakeHVACSystemGateway{
		state: testHVACSystem(t, domain.HVACSystemModeHeat),
	}
	service := NewHVACService(gateway)

	got, err := service.ApplyIntent(context.Background(), SetModeIntent{Mode: domain.HVACSystemMode("invalid")})

	require.ErrorIs(t, err, domain.ErrInvalidHVACMode)
	assert.Nil(t, got)
	assert.Equal(t, 1, gateway.getStateCalls)
	assert.Equal(t, 0, gateway.applyIntentCalls)
}

func TestHVACServiceApplyIntentReturnsGatewayError(t *testing.T) {
	wantErr := errors.New("apply failed")
	gateway := &fakeHVACSystemGateway{
		state:          testHVACSystem(t, domain.HVACSystemModeHeat),
		applyIntentErr: wantErr,
	}
	service := NewHVACService(gateway)

	got, err := service.ApplyIntent(context.Background(), SetModeIntent{Mode: domain.HVACSystemModeCool})

	require.ErrorIs(t, err, wantErr)
	assert.Nil(t, got)
	assert.Equal(t, 1, gateway.getStateCalls)
	assert.Equal(t, 1, gateway.applyIntentCalls)
}

func TestHVACServiceApplyIntentUsesExpectedResolvedIntents(t *testing.T) {
	testCases := []struct {
		name string
		act  func(*HVACService) (*domain.HVACSystem, error)
		want ResolvedIntent
	}{
		{
			name: "set room preset",
			act: func(service *HVACService) (*domain.HVACSystem, error) {
				return service.ApplyIntent(context.Background(), SetRoomPresetIntent{
					Room:   "Living Room",
					Preset: domain.PresetEco,
				})
			},
			want: SetRoomPresetIntent{
				Room:   "Living Room",
				Preset: domain.PresetEco,
			},
		},
		{
			name: "set room power on",
			act: func(service *HVACService) (*domain.HVACSystem, error) {
				return service.ApplyIntent(context.Background(), SetRoomPowerIntent{
					Room: "Living Room",
					On:   true,
				})
			},
			want: SetRoomPowerIntent{
				Room: "Living Room",
				On:   true,
			},
		},
		{
			name: "set room power off",
			act: func(service *HVACService) (*domain.HVACSystem, error) {
				return service.ApplyIntent(context.Background(), SetRoomPowerIntent{
					Room: "Living Room",
					On:   false,
				})
			},
			want: SetRoomPowerIntent{
				Room: "Living Room",
				On:   false,
			},
		},
		{
			name: "set temperature current/current",
			act: func(service *HVACService) (*domain.HVACSystem, error) {
				return service.ApplyIntent(context.Background(), SetTemperatureIntent{
					Room:    "Living Room",
					Preset:  TemperaturePresetCurrent,
					Mode:    TemperatureModeCurrent,
					Value:   21,
					IsDelta: false,
				})
			},
			want: ResolvedSetTemperatureIntent{
				Room:   "Living Room",
				Preset: domain.PresetComfort,
				Mode:   domain.HVACSystemModeHeat,
				Value:  21,
			},
		},
		{
			name: "set temperature explicit target",
			act: func(service *HVACService) (*domain.HVACSystem, error) {
				return service.ApplyIntent(context.Background(), SetTemperatureIntent{
					Room:    "Living Room",
					Preset:  TemperaturePresetEco,
					Mode:    TemperatureModeCool,
					Value:   27,
					IsDelta: false,
				})
			},
			want: ResolvedSetTemperatureIntent{
				Room:   "Living Room",
				Preset: domain.PresetEco,
				Mode:   domain.HVACSystemModeCool,
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
			assert.Equal(t, 1, gateway.applyIntentCalls)
			assert.Equal(t, tc.want, gateway.appliedIntent)
		})
	}
}

func TestHVACServiceApplyIntentTemperatureResolvesTargets(t *testing.T) {
	testCases := []struct {
		name string
		act  func(*HVACService) (*domain.HVACSystem, error)
		want ResolvedIntent
	}{
		{
			name: "set temperature current target becomes current setpoint change",
			act: func(service *HVACService) (*domain.HVACSystem, error) {
				return service.ApplyIntent(context.Background(), SetTemperatureIntent{
					Room:    "Living Room",
					Preset:  TemperaturePresetCurrent,
					Mode:    TemperatureModeCurrent,
					Value:   21,
					IsDelta: false,
				})
			},
			want: ResolvedSetTemperatureIntent{
				Room:   "Living Room",
				Preset: domain.PresetComfort,
				Mode:   domain.HVACSystemModeHeat,
				Value:  21,
			},
		},
		{
			name: "set temperature explicit target",
			act: func(service *HVACService) (*domain.HVACSystem, error) {
				return service.ApplyIntent(context.Background(), SetTemperatureIntent{
					Room:    "Living Room",
					Preset:  TemperaturePresetEco,
					Mode:    TemperatureModeCool,
					Value:   27,
					IsDelta: false,
				})
			},
			want: ResolvedSetTemperatureIntent{
				Room:   "Living Room",
				Preset: domain.PresetEco,
				Mode:   domain.HVACSystemModeCool,
				Value:  27,
			},
		},
		{
			name: "set temperature delta on resolved target",
			act: func(service *HVACService) (*domain.HVACSystem, error) {
				return service.ApplyIntent(context.Background(), SetTemperatureIntent{
					Room:    "Living Room",
					Preset:  TemperaturePresetComfort,
					Mode:    TemperatureModeHeat,
					Value:   1,
					IsDelta: true,
				})
			},
			want: ResolvedSetTemperatureIntent{
				Room:   "Living Room",
				Preset: domain.PresetComfort,
				Mode:   domain.HVACSystemModeHeat,
				Value:  21,
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
			assert.Equal(t, 1, gateway.applyIntentCalls)
			assert.Equal(t, tc.want, gateway.appliedIntent)
		})
	}
}
