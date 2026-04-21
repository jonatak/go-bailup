package domain_test

import (
	"testing"

	"github.com/jonatak/go-bailup/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHVACSystemValidatesMode(t *testing.T) {
	system, err := domain.NewHVACSystem(domain.HVACSystemModeHeat, []domain.Thermostat{
		mustThermostat(t, "Living Room", domain.PresetComfort),
	})

	require.NoError(t, err)
	assert.Equal(t, domain.HVACSystemModeHeat, system.Mode())
	assert.Len(t, system.Thermostats(), 1)
}

func TestNewHVACSystemRejectsInvalidMode(t *testing.T) {
	system, err := domain.NewHVACSystem(domain.HVACSystemMode("invalid"), nil)

	require.ErrorIs(t, err, domain.ErrInvalidHVACMode)
	assert.Nil(t, system)
}

func TestNewHVACSystemRejectsInvalidThermostat(t *testing.T) {
	system, err := domain.NewHVACSystem(domain.HVACSystemModeHeat, []domain.Thermostat{
		{},
	})

	require.ErrorIs(t, err, domain.ErrInvalidPresetMode)
	assert.Nil(t, system)
}

func TestHVACSystemValidateChecksModeAndThermostats(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	require.NoError(t, system.Validate())
}

func TestHVACSystemValidateRejectsInvalidZeroValueSystem(t *testing.T) {
	var system domain.HVACSystem

	require.ErrorIs(t, system.Validate(), domain.ErrInvalidHVACMode)
}

func TestNewHVACSystemCopiesThermostatSlice(t *testing.T) {
	livingRoom := mustThermostat(t, "Living Room", domain.PresetComfort)
	bedroom := mustThermostat(t, "Bedroom", domain.PresetEco)
	thermostats := []domain.Thermostat{livingRoom}
	system, err := domain.NewHVACSystem(domain.HVACSystemModeHeat, thermostats)
	require.NoError(t, err)

	thermostats[0] = bedroom

	current, err := system.CurrentSetpoint("Living Room")
	require.NoError(t, err)
	assert.Equal(t, 20.0, current)
}

func TestHVACSystemThermostatsReturnsCopy(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)
	thermostats := system.Thermostats()
	require.Len(t, thermostats, 1)

	thermostats[0] = mustThermostat(t, "Bedroom", domain.PresetEco)

	current, err := system.CurrentSetpoint("living room")
	require.NoError(t, err)
	assert.Equal(t, 20.0, current)
}

func TestHVACSystemRoomLookupIsCaseInsensitive(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	current, err := system.CurrentSetpoint("living room")

	require.NoError(t, err)
	assert.Equal(t, 20.0, current)
}

func TestHVACSystemRoomLookupReturnsErrorWhenRoomDoesNotExist(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	current, err := system.CurrentSetpoint("Kitchen")

	require.ErrorIs(t, err, domain.ErrThermostatNotFound)
	assert.Equal(t, 0.0, current)
}

func TestHVACSystemCurrentSetpointUsesCurrentModeAndPreset(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	current, err := system.CurrentSetpoint("Living Room")
	require.NoError(t, err)
	assert.Equal(t, 20.0, current)

	change, err := system.SetRoomPreset("Living Room", domain.PresetEco)
	require.NoError(t, err)
	assert.Equal(t, domain.RoomPresetChanged{
		Room:   "Living Room",
		Preset: domain.PresetEco,
	}, change)
	assert.Equal(t, domain.ChangeRoomPreset, change.Kind())

	current, err = system.CurrentSetpoint("Living Room")
	require.NoError(t, err)
	assert.Equal(t, 18.0, current)

	change, err = system.SetMode(domain.HVACSystemModeCool)
	require.NoError(t, err)
	assert.Equal(t, domain.HVACModeChanged{
		Mode: domain.HVACSystemModeCool,
	}, change)
	assert.Equal(t, domain.ChangeHVACMode, change.Kind())

	current, err = system.CurrentSetpoint("Living Room")
	require.NoError(t, err)
	assert.Equal(t, 26.0, current)
}

func TestHVACSystemCurrentSetpointRejectsModeWithoutSetpoints(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeOff)

	current, err := system.CurrentSetpoint("Living Room")

	require.ErrorIs(t, err, domain.ErrCurrentSetpointUnavailable)
	assert.Equal(t, 0.0, current)
}

func TestHVACSystemSetpointUsesTargetModeAndPreset(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	testCases := []struct {
		name   string
		mode   domain.HVACSystemMode
		preset domain.ThermostatPreset
		want   float64
	}{
		{
			name:   "heat comfort",
			mode:   domain.HVACSystemModeHeat,
			preset: domain.PresetComfort,
			want:   20,
		},
		{
			name:   "heat eco",
			mode:   domain.HVACSystemModeHeat,
			preset: domain.PresetEco,
			want:   18,
		},
		{
			name:   "cool comfort",
			mode:   domain.HVACSystemModeCool,
			preset: domain.PresetComfort,
			want:   24,
		},
		{
			name:   "cool eco",
			mode:   domain.HVACSystemModeCool,
			preset: domain.PresetEco,
			want:   26,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := system.Setpoint("Living Room", tc.mode, tc.preset)

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestHVACSystemSetpointRejectsModeWithoutSetpoints(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	current, err := system.Setpoint("Living Room", domain.HVACSystemModeOff, domain.PresetComfort)

	require.ErrorIs(t, err, domain.ErrInvalidTemperatureSettingForMode)
	assert.Equal(t, 0.0, current)
}

func TestHVACSystemSetModeDoesNotChangeStateForInvalidMode(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	change, err := system.SetMode(domain.HVACSystemMode("invalid"))

	require.ErrorIs(t, err, domain.ErrInvalidHVACMode)
	assert.Nil(t, change)
	assert.Equal(t, domain.HVACSystemModeHeat, system.Mode())
}

func TestHVACSystemSetCurrentSetpointUpdatesActivePreset(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	change, err := system.SetCurrentSetpoint("Living Room", 21)
	require.NoError(t, err)
	assert.Equal(t, domain.TemperatureChanged{
		Room:   "Living Room",
		Mode:   domain.HVACSystemModeHeat,
		Preset: domain.PresetComfort,
		Value:  21,
	}, change)
	assert.Equal(t, domain.ChangeTemperature, change.Kind())

	current, err := system.CurrentSetpoint("Living Room")
	require.NoError(t, err)
	assert.Equal(t, 21.0, current)
}

func TestHVACSystemSetTemperatureRejectsInvalidComfortEcoRange(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	change, err := system.SetTemperature("Living Room", domain.HVACSystemModeHeat, domain.PresetEco, 19)

	require.ErrorIs(t, err, domain.ErrSetpointUnsupportedForMode)
	assert.Nil(t, change)

	current, currentErr := system.CurrentSetpoint("Living Room")
	require.NoError(t, currentErr)
	assert.Equal(t, 20.0, current)
}

func TestHVACSystemSetTemperatureUpdatesTargetModeAndPreset(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	change, err := system.SetTemperature("Living Room", domain.HVACSystemModeCool, domain.PresetEco, 27)
	require.NoError(t, err)
	assert.Equal(t, domain.TemperatureChanged{
		Room:   "Living Room",
		Mode:   domain.HVACSystemModeCool,
		Preset: domain.PresetEco,
		Value:  27,
	}, change)
	assert.Equal(t, domain.ChangeTemperature, change.Kind())

	_, err = system.SetMode(domain.HVACSystemModeCool)
	require.NoError(t, err)
	_, err = system.SetRoomPreset("Living Room", domain.PresetEco)
	require.NoError(t, err)

	current, err := system.CurrentSetpoint("Living Room")
	require.NoError(t, err)
	assert.Equal(t, 27.0, current)
}

func TestHVACSystemSetRoomPresetValidatesPreset(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	change, err := system.SetRoomPreset("Living Room", domain.ThermostatPreset("away"))

	require.ErrorIs(t, err, domain.ErrInvalidPresetMode)
	assert.Nil(t, change)
	current, currentErr := system.CurrentSetpoint("Living Room")
	require.NoError(t, currentErr)
	assert.Equal(t, 20.0, current)
}

func TestHVACSystemTurnRoomOnAndOff(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	change, err := system.TurnRoomOff("Living Room")
	require.NoError(t, err)
	assert.Equal(t, domain.RoomPowerChanged{
		Room: "Living Room",
		On:   false,
	}, change)
	assert.Equal(t, domain.ChangeRoomPower, change.Kind())
	thermostats := system.Thermostats()
	require.Len(t, thermostats, 1)
	assert.False(t, thermostats[0].IsOn())

	change, err = system.TurnRoomOn("Living Room")
	require.NoError(t, err)
	assert.Equal(t, domain.RoomPowerChanged{
		Room: "Living Room",
		On:   true,
	}, change)
	assert.Equal(t, domain.ChangeRoomPower, change.Kind())
	thermostats = system.Thermostats()
	require.Len(t, thermostats, 1)
	assert.True(t, thermostats[0].IsOn())
}

func mustHVACSystem(t *testing.T, mode domain.HVACSystemMode) *domain.HVACSystem {
	t.Helper()

	system, err := domain.NewHVACSystem(mode, []domain.Thermostat{
		mustThermostat(t, "Living Room", domain.PresetComfort),
	})
	require.NoError(t, err)

	return system
}
