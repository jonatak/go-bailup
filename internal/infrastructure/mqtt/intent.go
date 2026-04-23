package mqtt

import (
	"github.com/jonatak/go-bailup/internal/application"
	"github.com/jonatak/go-bailup/internal/domain"
)

type intent interface {
	Apply(*application.HVACService) (*domain.HVACSystem, error)
}

type setTemperatureIntent struct {
	room  string
	value float64
}

type setModeIntent struct {
	mode domain.HVACSystemMode
}

type turnRoomOnIntent struct {
	room string
}

type turnRoomOffIntent struct {
	room string
}

type setPresetIntent struct {
	room   string
	preset domain.ThermostatPreset
}

func (i setTemperatureIntent) Apply(service *application.HVACService) (*domain.HVACSystem, error) {
	return service.SetCurrentSetpoint(i.room, i.value)
}

func (i setModeIntent) Apply(service *application.HVACService) (*domain.HVACSystem, error) {
	return service.SetMode(i.mode)
}

func (i turnRoomOnIntent) Apply(service *application.HVACService) (*domain.HVACSystem, error) {
	return service.TurnRoomOn(i.room)
}

func (i turnRoomOffIntent) Apply(service *application.HVACService) (*domain.HVACSystem, error) {
	return service.TurnRoomOff(i.room)
}

func (i setPresetIntent) Apply(service *application.HVACService) (*domain.HVACSystem, error) {
	return service.SetRoomPreset(i.room, i.preset)
}
