package bailup

import (
	"testing"

	"github.com/jonatak/go-bailup/internal/domain"
	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/command"
	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHVACSystemFromState(t *testing.T) {
	state := mapperTestState()

	system, err := HVACSystemFromState(state)

	require.NoError(t, err)
	assert.Equal(t, domain.HVACSystemModeHeat, system.Mode())

	thermostats := system.Thermostats()
	require.Len(t, thermostats, 2)

	assert.Equal(t, "Living Room", thermostats[0].Room())
	assert.Equal(t, domain.PresetComfort, thermostats[0].Preset())
	assert.True(t, thermostats[0].IsOn())
	assert.False(t, thermostats[0].IsRunning())
	assert.Equal(t, 20.0, thermostats[0].HeatSetting().Comfort())
	assert.Equal(t, 18.0, thermostats[0].HeatSetting().Eco())
	assert.Equal(t, 24.0, thermostats[0].CoolSetting().Comfort())
	assert.Equal(t, 26.0, thermostats[0].CoolSetting().Eco())

	assert.Equal(t, "Bedroom", thermostats[1].Room())
	assert.Equal(t, domain.PresetEco, thermostats[1].Preset())
	assert.False(t, thermostats[1].IsOn())
	assert.True(t, thermostats[1].IsRunning())
}

func TestHVACSystemFromStateRejectsNilState(t *testing.T) {
	system, err := HVACSystemFromState(nil)

	require.Error(t, err)
	assert.Nil(t, system)
	assert.Contains(t, err.Error(), "state is nil")
}

func TestCommandFromHVACModeChanged(t *testing.T) {
	cmd, err := CommandFromChange(mapperTestState(), domain.HVACModeChanged{
		Mode: domain.HVACSystemModeCool,
	})

	require.NoError(t, err)
	assert.Equal(t, command.ModeCommand{
		Value: int(model.UCModeCool),
	}, cmd)
}

func TestCommandFromRoomPresetChanged(t *testing.T) {
	cmd, err := CommandFromChange(mapperTestState(), domain.RoomPresetChanged{
		Room:   "living room",
		Preset: domain.PresetEco,
	})

	require.NoError(t, err)
	assert.Equal(t, command.Preset{
		ThermostatID: 9152,
		Value:        int(model.ThModeEco),
	}, cmd)
}

func TestCommandFromRoomPowerChanged(t *testing.T) {
	cmd, err := CommandFromChange(mapperTestState(), domain.RoomPowerChanged{
		Room: "bedroom",
		On:   true,
	})

	require.NoError(t, err)
	assert.Equal(t, command.RoomPowerCommand{
		ThermostatID: 9154,
		On:           true,
	}, cmd)
}

func TestCommandFromTemperatureChanged(t *testing.T) {
	cmd, err := CommandFromChange(mapperTestState(), domain.TemperatureChanged{
		Room:   "living room",
		Mode:   domain.HVACSystemModeCool,
		Preset: domain.PresetEco,
		Value:  26.5,
	})

	require.NoError(t, err)
	assert.Equal(t, command.TemperatureCommand{
		ThermostatID: 9152,
		UCMode:       model.UCModeCool,
		ThMode:       model.ThModeEco,
		Value:        26.5,
	}, cmd)
}

func TestCommandFromChangeRejectsNilState(t *testing.T) {
	cmd, err := CommandFromChange(nil, domain.HVACModeChanged{
		Mode: domain.HVACSystemModeCool,
	})

	require.Error(t, err)
	assert.Nil(t, cmd)
	assert.Contains(t, err.Error(), "state is nil")
}

func TestCommandFromChangeRejectsNilChange(t *testing.T) {
	cmd, err := CommandFromChange(mapperTestState(), nil)

	require.Error(t, err)
	assert.Nil(t, cmd)
	assert.Contains(t, err.Error(), "change is nil")
}

func TestCommandFromTemperatureChangedRejectsUnknownRoom(t *testing.T) {
	cmd, err := CommandFromChange(mapperTestState(), domain.TemperatureChanged{
		Room:   "kitchen",
		Mode:   domain.HVACSystemModeCool,
		Preset: domain.PresetEco,
		Value:  26.5,
	})

	require.Error(t, err)
	assert.Nil(t, cmd)
	assert.Contains(t, err.Error(), "thermostat \"kitchen\" not found")
}

func mapperTestState() *model.State {
	return &model.State{
		ID:          2890,
		UCMode:      model.UCModeHeat,
		IsConnected: true,
		Thermostats: []model.Thermostat{
			{
				ID:             9152,
				Key:            "th1",
				Number:         1,
				Name:           "Living Room",
				Temperature:    20.4,
				IsOn:           true,
				SetpointHotT1:  20,
				SetpointHotT2:  18,
				SetpointCoolT1: 24,
				SetpointCoolT2: 26,
				MotorState:     0,
				T1T2:           model.ThModeComfort,
				IsConnected:    true,
			},
			{
				ID:             9154,
				Key:            "th2",
				Number:         2,
				Name:           "Bedroom",
				Temperature:    19.2,
				IsOn:           false,
				SetpointHotT1:  19,
				SetpointHotT2:  17,
				SetpointCoolT1: 25,
				SetpointCoolT2: 27,
				MotorState:     5,
				T1T2:           model.ThModeEco,
				IsConnected:    true,
			},
		},
	}
}
