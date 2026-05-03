package mqtt

import (
	"fmt"

	"github.com/jonatak/baillconnect-to-mqtt/internal/domain"
)

type MQTTBatterySensor struct {
	DefaultEntityID   string `json:"default_entity_id"`
	DeviceClass       string `json:"device_class"`
	EntityCategory    string `json:"entity_category"`
	StateClass        string `json:"state_class"`
	StateTopic        string `json:"state_topic"`
	UniqueID          string `json:"unique_id"`
	UnitOfMeasurement string `json:"unit_of_measurement"`
	Device            Device `json:"device"`
}

func BatteryFromThermostatDomain(t domain.Thermostat, prefix string) MQTTBatterySensor {
	return MQTTBatterySensor{
		DefaultEntityID:   fmt.Sprintf("sensor.thermostat_%s_battery", slugify(t.Room())),
		DeviceClass:       "battery",
		EntityCategory:    "diagnostic",
		StateClass:        "measurement",
		StateTopic:        fmt.Sprintf("%s/th_%d/battery", prefix, t.ID()),
		UniqueID:          fmt.Sprintf("%s-%s-%d-battery", prefix, slugify(t.Room()), t.ID()),
		UnitOfMeasurement: "%",
		Device: Device{
			Identifiers:   []string{fmt.Sprintf("bailup_%d", t.ID())},
			Manufacturer:  "Bail Industry",
			Name:          fmt.Sprintf("Thermostat %s", t.Room()),
			SuggestedArea: t.Room(),
		},
	}
}
