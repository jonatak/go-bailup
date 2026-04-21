package domain

type TemperatureSettings struct {
	comfort float64
	eco     float64
}

func NewTemperatureSettings(mode HVACSystemMode, comfort float64, eco float64) (TemperatureSettings, error) {
	settings := TemperatureSettings{
		comfort: comfort,
		eco:     eco,
	}

	if err := settings.validateForMode(mode); err != nil {
		return TemperatureSettings{}, err
	}

	return settings, nil
}

func (t TemperatureSettings) Comfort() float64 {
	return t.comfort
}

func (t TemperatureSettings) Eco() float64 {
	return t.eco
}

func (t TemperatureSettings) validateForMode(mode HVACSystemMode) error {
	switch mode {
	case HVACSystemModeCool:
		if t.eco < t.comfort {
			return ErrEcoMustBeBiggerThanComfort
		}
		if (t.eco - t.comfort) < 2 {
			return ErrSetpointUnsupportedForMode
		}
	case HVACSystemModeHeat:
		if t.eco > t.comfort {
			return ErrComfortMustBeBiggerThanEco
		}
		if (t.comfort - t.eco) < 2 {
			return ErrSetpointUnsupportedForMode
		}
	default:
		return ErrSetpointUnsupportedForMode
	}
	return nil
}
