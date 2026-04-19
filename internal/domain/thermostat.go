package domain

type Thermostat struct {
	room        string
	preset      ThermostatPreset
	isOn        bool
	isRunning   bool
	coolSetting TemperatureSettings
	heatSetting TemperatureSettings
}

func NewThermostat(
	room string,
	preset ThermostatPreset,
	isOn bool,
	isRunning bool,
	heatSetting TemperatureSettings,
	coolSetting TemperatureSettings,
) (Thermostat, error) {
	thermostat := Thermostat{
		room:        room,
		preset:      preset,
		isOn:        isOn,
		isRunning:   isRunning,
		heatSetting: heatSetting,
		coolSetting: coolSetting,
	}

	if err := thermostat.Validate(); err != nil {
		return Thermostat{}, err
	}

	return thermostat, nil
}

func (t Thermostat) Validate() error {
	if err := t.preset.Validate(); err != nil {
		return err
	}

	if err := t.heatSetting.validateForMode(HVACSystemModeHeat); err != nil {
		return err
	}

	if err := t.coolSetting.validateForMode(HVACSystemModeCool); err != nil {
		return err
	}

	return nil
}

func (t *Thermostat) Room() string {
	return t.room
}

func (t *Thermostat) Preset() ThermostatPreset {
	return t.preset
}

func (t *Thermostat) IsOn() bool {
	return t.isOn
}

func (t *Thermostat) IsRunning() bool {
	return t.isRunning
}

func (t *Thermostat) CoolSetting() TemperatureSettings {
	return t.coolSetting
}

func (t *Thermostat) HeatSetting() TemperatureSettings {
	return t.heatSetting
}

func (t *Thermostat) setPreset(preset ThermostatPreset) error {
	if err := preset.Validate(); err != nil {
		return err
	}
	t.preset = preset
	return nil
}

func (t *Thermostat) turnOn() {
	t.isOn = true
}

func (t *Thermostat) turnOff() {
	t.isOn = false
}

func (t *Thermostat) currentSetpointForMode(mode HVACSystemMode) (float64, error) {
	if err := t.preset.Validate(); err != nil {
		return 0, err
	}

	return t.setpointFor(mode, t.preset)
}

func (t *Thermostat) setpointFor(mode HVACSystemMode, preset ThermostatPreset) (float64, error) {
	if err := preset.Validate(); err != nil {
		return 0, err
	}

	switch mode {
	case HVACSystemModeHeat:
		if preset == PresetEco {
			return t.heatSetting.eco, nil
		}
		return t.heatSetting.comfort, nil

	case HVACSystemModeCool:
		if preset == PresetEco {
			return t.coolSetting.eco, nil
		}
		return t.coolSetting.comfort, nil

	default:
		return 0, ErrInvalidTemperatureSettingForMode
	}
}

func (t *Thermostat) currentSettingsForMode(mode HVACSystemMode) (TemperatureSettings, error) {
	switch mode {
	case HVACSystemModeCool:
		return t.coolSetting, nil
	case HVACSystemModeHeat:
		return t.heatSetting, nil
	default:
		return TemperatureSettings{}, ErrInvalidTemperatureSettingForMode
	}
}

func (t *Thermostat) setTemperatureSettingForMode(mode HVACSystemMode, settings TemperatureSettings) error {
	if err := settings.validateForMode(mode); err != nil {
		return err
	}
	if mode == HVACSystemModeCool {
		t.coolSetting = settings
		return nil
	}
	if mode == HVACSystemModeHeat {
		t.heatSetting = settings
		return nil
	}

	return ErrInvalidTemperatureSettingForMode
}

func (t *Thermostat) setTemperature(mode HVACSystemMode, preset ThermostatPreset, temp float64) error {
	if err := mode.Validate(); err != nil {
		return err
	}

	if err := preset.Validate(); err != nil {
		return err
	}

	if mode != HVACSystemModeCool && mode != HVACSystemModeHeat {
		return ErrInvalidTemperatureSettingForMode
	}

	tempSettings, err := t.currentSettingsForMode(mode)
	if err != nil {
		return err
	}

	if preset == PresetComfort {
		tempSettings.comfort = temp
	} else {
		tempSettings.eco = temp
	}

	return t.setTemperatureSettingForMode(mode, tempSettings)
}
