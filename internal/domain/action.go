package domain

type ThermostatAction string

const (
	ThermostatActionOff     ThermostatAction = "off"
	ThermostatActionIdle    ThermostatAction = "idle"
	ThermostatActionCooling ThermostatAction = "cooling"
	ThermostatActionHeating ThermostatAction = "heating"
	ThermostatActionDrying  ThermostatAction = "drying"
	ThermostatActionFan     ThermostatAction = "fan"
)
