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

	found, err := system.FindThermostat("Living Room")
	require.NoError(t, err)
	assert.Equal(t, "Living Room", found.Room())
}

func TestHVACSystemThermostatsReturnsCopy(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)
	thermostats := system.Thermostats()
	require.Len(t, thermostats, 1)

	require.NoError(t, thermostats[0].SetPreset(domain.PresetEco))

	current, err := system.CurrentSetpoint("living room")
	require.NoError(t, err)
	assert.Equal(t, 20.0, current)
}

func TestHVACSystemFindThermostatIsCaseInsensitive(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	thermostat, err := system.FindThermostat("living room")

	require.NoError(t, err)
	assert.Equal(t, "Living Room", thermostat.Room())
}

func TestHVACSystemFindThermostatReturnsErrorWhenRoomDoesNotExist(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	thermostat, err := system.FindThermostat("Kitchen")

	require.ErrorIs(t, err, domain.ErrThermostatNotFound)
	assert.Nil(t, thermostat)
}

func TestHVACSystemCurrentSetpointUsesCurrentModeAndPreset(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	current, err := system.CurrentSetpoint("Living Room")
	require.NoError(t, err)
	assert.Equal(t, 20.0, current)

	thermostat, err := system.FindThermostat("Living Room")
	require.NoError(t, err)
	require.NoError(t, thermostat.SetPreset(domain.PresetEco))

	current, err = system.CurrentSetpoint("Living Room")
	require.NoError(t, err)
	assert.Equal(t, 18.0, current)

	require.NoError(t, system.SetMode(domain.HVACSystemModeCool))

	current, err = system.CurrentSetpoint("Living Room")
	require.NoError(t, err)
	assert.Equal(t, 26.0, current)
}

func TestHVACSystemCurrentSetpointRejectsModeWithoutSetpoints(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeOff)

	current, err := system.CurrentSetpoint("Living Room")

	require.ErrorIs(t, err, domain.ErrCurrentSetPointInvalid)
	assert.Equal(t, 0.0, current)
}

func TestHVACSystemSetModeDoesNotChangeStateForInvalidMode(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	err := system.SetMode(domain.HVACSystemMode("invalid"))

	require.ErrorIs(t, err, domain.ErrInvalidHVACMode)
	assert.Equal(t, domain.HVACSystemModeHeat, system.Mode())
}

func TestHVACSystemSetCurrentSetPointUpdatesActivePreset(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	require.NoError(t, system.SetCurrentSetPoint("Living Room", 21))

	current, err := system.CurrentSetpoint("Living Room")
	require.NoError(t, err)
	assert.Equal(t, 21.0, current)
}

func TestHVACSystemSetTemperatureRejectsInvalidComfortEcoRange(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	err := system.SetHeatEcoTemp("Living Room", 19)

	require.ErrorIs(t, err, domain.ErrInvalidTemperatureRange)

	current, currentErr := system.CurrentSetpoint("Living Room")
	require.NoError(t, currentErr)
	assert.Equal(t, 20.0, current)
}

func TestHVACSystemSetTemperatureUpdatesTargetModeAndPreset(t *testing.T) {
	system := mustHVACSystem(t, domain.HVACSystemModeHeat)

	require.NoError(t, system.SetCoolEcoTemp("Living Room", 27))
	require.NoError(t, system.SetMode(domain.HVACSystemModeCool))
	thermostat, err := system.FindThermostat("Living Room")
	require.NoError(t, err)
	require.NoError(t, thermostat.SetPreset(domain.PresetEco))

	current, err := system.CurrentSetpoint("Living Room")
	require.NoError(t, err)
	assert.Equal(t, 27.0, current)
}

func mustHVACSystem(t *testing.T, mode domain.HVACSystemMode) *domain.HVACSystem {
	t.Helper()

	system, err := domain.NewHVACSystem(mode, []domain.Thermostat{
		mustThermostat(t, "Living Room", domain.PresetComfort),
	})
	require.NoError(t, err)

	return system
}
