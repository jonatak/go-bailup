package command

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/jonatak/go-bailup/internal/bailup/model"
)

type TemperatureCommand struct {
	ThermostatID int
	Value        float64
	UCMode       model.UCMode
	ThMode       model.ThMode
}

func (p TemperatureCommand) ToJSON() ([]byte, error) {
	value := math.Round(p.Value*10) / 10

	payload := map[string]float64{}
	switch {
	case p.UCMode == model.UCModeHeat && p.ThMode == model.ThModeComfort:
		payload[fmt.Sprintf("thermostats.%d.setpoint_hot_t1", p.ThermostatID)] = value
	case p.UCMode == model.UCModeHeat && p.ThMode == model.ThModeEco:
		payload[fmt.Sprintf("thermostats.%d.setpoint_hot_t2", p.ThermostatID)] = value
	case p.UCMode == model.UCModeCool && p.ThMode == model.ThModeComfort:
		payload[fmt.Sprintf("thermostats.%d.setpoint_cool_t1", p.ThermostatID)] = value
	case p.UCMode == model.UCModeCool && p.ThMode == model.ThModeEco:
		payload[fmt.Sprintf("thermostats.%d.setpoint_cool_t2", p.ThermostatID)] = value
	default:
		return nil, fmt.Errorf("unsupported temperature combination: uc_mode=%s th_mode=%s", p.UCMode, p.ThMode)
	}

	return json.Marshal(payload)
}
