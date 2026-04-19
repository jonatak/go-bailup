package domain

type ChangeKind string

const (
	ChangeHVACMode    ChangeKind = "hvac_mode"
	ChangeRoomPreset  ChangeKind = "room_preset"
	ChangeRoomPower   ChangeKind = "room_power"
	ChangeTemperature ChangeKind = "temperature"
)

type Change interface {
	Kind() ChangeKind
}

type HVACModeChanged struct {
	Mode HVACSystemMode
}

type RoomPresetChanged struct {
	Room   string
	Preset ThermostatPreset
}

type RoomPowerChanged struct {
	Room string
	On   bool
}

type TemperatureChanged struct {
	Room   string
	Mode   HVACSystemMode
	Preset ThermostatPreset
	Value  float64
}

func (HVACModeChanged) Kind() ChangeKind {
	return ChangeHVACMode
}

func (RoomPresetChanged) Kind() ChangeKind {
	return ChangeRoomPreset
}

func (RoomPowerChanged) Kind() ChangeKind {
	return ChangeRoomPower
}

func (TemperatureChanged) Kind() ChangeKind {
	return ChangeTemperature
}
