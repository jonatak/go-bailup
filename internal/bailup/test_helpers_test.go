package bailup_test

import "github.com/jonatak/go-bailup/internal/bailup/model"

func testState() *model.State {
	return &model.State{
		ID:          2890,
		UCMode:      model.UCModeHeat,
		IsConnected: true,
		Thermostats: []model.Thermostat{
			{
				ID:             9152,
				Key:            "th1",
				Number:         1,
				Name:           "Living Room",
				Temperature:    20.4,
				Zone:           1,
				IsOn:           true,
				SetpointHotT1:  20,
				SetpointHotT2:  18,
				SetpointCoolT1: 24,
				SetpointCoolT2: 26,
				MotorState:     0,
				T1T2:           model.ThModeComfort,
				IsBatteryLow:   false,
				IsConnected:    true,
			},
			{
				ID:             9154,
				Key:            "th2",
				Number:         2,
				Name:           "Bedroom",
				Temperature:    19.2,
				Zone:           1,
				IsOn:           false,
				SetpointHotT1:  19,
				SetpointHotT2:  17,
				SetpointCoolT1: 25,
				SetpointCoolT2: 27,
				MotorState:     0,
				T1T2:           model.ThModeEco,
				IsBatteryLow:   false,
				IsConnected:    true,
			},
		},
	}
}
