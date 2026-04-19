package domain

type HVACSystemMode string

const (
	HVACSystemModeOff     HVACSystemMode = "off"
	HVACSystemModeCool    HVACSystemMode = "cool"
	HVACSystemModeHeat    HVACSystemMode = "heat"
	HVACSystemModeDry     HVACSystemMode = "dry"
	HVACSystemModeFanOnly HVACSystemMode = "fan-only"
)

func (h HVACSystemMode) Validate() error {
	switch h {
	case HVACSystemModeOff, HVACSystemModeCool, HVACSystemModeHeat, HVACSystemModeDry, HVACSystemModeFanOnly:
		return nil
	default:
		return ErrInvalidHVACMode
	}
}
