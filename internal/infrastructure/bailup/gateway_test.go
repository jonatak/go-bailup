package bailup

import (
	"context"
	"errors"
	"testing"

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

	system, err := gateway.GetHVACSystemState(context.Background())

	require.Error(t, err)
	assert.Nil(t, system)
	assert.True(t, errors.Is(err, application.ErrStateUnavailable))
}

func TestGatewayGetHVACSystemStateWrapsStateMappingError(t *testing.T) {
	gateway := &Gateway{
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

	system, err := gateway.ApplyResolvedIntent(context.Background(), application.SetModeIntent{
		Mode: domain.HVACSystemModeCool,
	})

	require.Error(t, err)
	assert.Nil(t, system)
	assert.True(t, errors.Is(err, application.ErrStateUnavailable))
}

func TestGatewayApplyResolvedIntentWrapsIntentMappingError(t *testing.T) {
	gateway := &Gateway{
		state: mapperTestState(),
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
