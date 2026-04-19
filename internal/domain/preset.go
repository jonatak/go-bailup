package domain

type ThermostatPreset string

const (
	PresetComfort ThermostatPreset = "comfort"
	PresetEco     ThermostatPreset = "eco"
)

func (p ThermostatPreset) Validate() error {
	switch p {
	case PresetComfort, PresetEco:
		return nil
	default:
		return ErrInvalidPresetMode
	}
}
