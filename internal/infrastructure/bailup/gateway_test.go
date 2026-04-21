package bailup

import (
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

	system, err := gateway.GetHVACSystemState()

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

	system, err := gateway.GetHVACSystemState()

	require.Error(t, err)
	assert.Nil(t, system)
	assert.True(t, errors.Is(err, application.ErrStateUnavailable))
	assert.True(t, errors.Is(err, domain.ErrComfortMustBeBiggerThanEco))
}

func TestGatewayApplyChangeWrapsStateLoadError(t *testing.T) {
	gateway := &Gateway{
		client: NewBailup("", "", ""),
	}

	system, err := gateway.ApplyChange(domain.HVACModeChanged{
		Mode: domain.HVACSystemModeCool,
	})

	require.Error(t, err)
	assert.Nil(t, system)
	assert.True(t, errors.Is(err, application.ErrStateUnavailable))
}

func TestGatewayApplyChangeWrapsChangeMappingError(t *testing.T) {
	gateway := &Gateway{
		state: mapperTestState(),
	}

	system, err := gateway.ApplyChange(domain.RoomPowerChanged{
		Room: "Kitchen",
		On:   true,
	})

	require.Error(t, err)
	assert.Nil(t, system)
	assert.True(t, errors.Is(err, application.ErrChangeRejected))
	assert.Contains(t, err.Error(), `thermostat "Kitchen" not found`)
}
