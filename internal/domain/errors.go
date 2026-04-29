package domain

import "errors"

var (
	ErrInvalidHVACMode                  = errors.New("invalid HVAC mode")
	ErrInvalidPresetMode                = errors.New("invalid Preset mode")
	ErrThermostatNotFound               = errors.New("thermostat not found")
	ErrCurrentSetpointUnavailable       = errors.New("current setpoint is only available in heat or cool mode")
	ErrInvalidTemperatureSettingForMode = errors.New("temperature setpoints are only available in heat or cool mode")
	ErrEcoMustBeBiggerThanComfort       = errors.New("cool eco setpoint must be greater than or equal to cool comfort setpoint")
	ErrComfortMustBeBiggerThanEco       = errors.New("heat comfort setpoint must be greater than or equal to heat eco setpoint")
	ErrSetpointUnsupportedForMode       = errors.New("comfort and eco setpoints must differ by at least 2 degrees")
)
