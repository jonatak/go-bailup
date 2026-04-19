package model

type Response struct {
	Data State `json:"data"`
}

type Thermostat struct {
	ID             int     `json:"id"`
	Key            string  `json:"key"`
	Number         int     `json:"number"`
	Name           string  `json:"name"`
	Temperature    float64 `json:"temperature"`
	Zone           int     `json:"zone"`
	IsOn           bool    `json:"is_on"`
	SetpointHotT1  float64 `json:"setpoint_hot_t1"`
	SetpointHotT2  float64 `json:"setpoint_hot_t2"`
	SetpointCoolT1 float64 `json:"setpoint_cool_t1"`
	SetpointCoolT2 float64 `json:"setpoint_cool_t2"`
	MotorState     int     `json:"motor_state"`
	T1T2           ThMode  `json:"t1_t2"`
	IsBatteryLow   bool    `json:"is_battery_low"`
	IsConnected    bool    `json:"is_connected"`
}
