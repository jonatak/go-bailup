package bailup

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jonatak/go-bailup/internal/application"
	"github.com/jonatak/go-bailup/internal/domain"
	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGatewayGetHVACSystemStateWrapsStateLoadError(t *testing.T) {
	gateway := &Gateway{
		client: NewBailup("", "", ""),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	system, err := gateway.GetHVACSystemState(ctx)

	require.Error(t, err)
	assert.Nil(t, system)
	assert.True(t, errors.Is(err, application.ErrStateUnavailable))
}

func TestGatewayGetHVACSystemStateWrapsStateMappingError(t *testing.T) {
	gateway := &Gateway{
		lastRefreshed: time.Now(),
		state: &model.State{
			UCMode: model.UCModeHeat,
			Thermostats: []model.Thermostat{
				{
					Name:           "Living Room",
					T1T2:           model.ThModeComfort,
					SetpointHotT1:  18,
					SetpointHotT2:  20,
					SetpointCoolT1: 24,
					SetpointCoolT2: 26,
				},
			},
		},
	}

	system, err := gateway.GetHVACSystemState(context.Background())

	require.Error(t, err)
	assert.Nil(t, system)
	assert.True(t, errors.Is(err, application.ErrStateUnavailable))
	assert.True(t, errors.Is(err, domain.ErrComfortMustBeBiggerThanEco))
}

func TestGatewayApplyResolvedIntentWrapsStateLoadError(t *testing.T) {
	gateway := &Gateway{
		client: NewBailup("", "", ""),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	system, err := gateway.ApplyResolvedIntent(ctx, application.SetModeIntent{
		Mode: domain.HVACSystemModeCool,
	})

	require.Error(t, err)
	assert.Nil(t, system)
	assert.True(t, errors.Is(err, application.ErrStateUnavailable))
}

func TestGatewayApplyResolvedIntentWrapsIntentMappingError(t *testing.T) {
	gateway := &Gateway{
		state:         mapperTestState(),
		lastRefreshed: time.Now(),
	}

	system, err := gateway.ApplyResolvedIntent(context.Background(), application.SetRoomPowerIntent{
		Room: "Kitchen",
		On:   true,
	})

	require.Error(t, err)
	assert.Nil(t, system)
	assert.True(t, errors.Is(err, application.ErrChangeRejected))
	assert.Contains(t, err.Error(), `thermostat "Kitchen" not found`)
}
