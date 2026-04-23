package application

import "github.com/jonatak/go-bailup/internal/domain"

type Intent interface {
	isIntent()
}

type ResolvedIntent interface {
	Intent
	isResolvedIntent()
}

type TemperatureModeTarget string

const (
	TemperatureModeCurrent TemperatureModeTarget = "current"
	TemperatureModeCool    TemperatureModeTarget = "cool"
	TemperatureModeHeat    TemperatureModeTarget = "heat"
)

type TemperaturePresetTarget string

const (
	TemperaturePresetCurrent TemperaturePresetTarget = "current"
	TemperaturePresetComfort TemperaturePresetTarget = "comfort"
	TemperaturePresetEco     TemperaturePresetTarget = "eco"
)

type SetTemperatureIntent struct {
	Room    string
	Preset  TemperaturePresetTarget
	Mode    TemperatureModeTarget
	Value   float64
	IsDelta bool
}

type SetModeIntent struct {
	Mode domain.HVACSystemMode
}

type SetRoomPowerIntent struct {
	Room string
	On   bool
}

type SetRoomPresetIntent struct {
	Room   string
	Preset domain.ThermostatPreset
}

type ResolvedSetTemperatureIntent struct {
	Room   string
	Preset domain.ThermostatPreset
	Mode   domain.HVACSystemMode
	Value  float64
}

func (SetTemperatureIntent) isIntent() {}
func (SetModeIntent) isIntent()        {}
func (SetRoomPowerIntent) isIntent()   {}
func (SetRoomPresetIntent) isIntent()  {}

func (SetModeIntent) isResolvedIntent()                {}
func (SetRoomPowerIntent) isResolvedIntent()           {}
func (SetRoomPresetIntent) isResolvedIntent()          {}
func (ResolvedSetTemperatureIntent) isIntent()         {}
func (ResolvedSetTemperatureIntent) isResolvedIntent() {}
