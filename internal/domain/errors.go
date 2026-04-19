package domain

import "errors"

var (
	ErrInvalidHVACMode                  = errors.New("Invalid HVAC mode")
	ErrInvalidPresetMode                = errors.New("Invalid Preset mode")
	ErrThermostatNotFound               = errors.New("Thermostat not found")
	ErrCurrentSetPointInvalid           = errors.New("HVAC system need to be in heat or cool mode for thermostat to have a current set point")
	ErrInvalidTemperatureSettingForMode = errors.New("Can't set temperature for the targetted mode")
	ErrInvalidTemperatureRange          = errors.New("Different between comfort and eco should be at least 2 degrees")
)
