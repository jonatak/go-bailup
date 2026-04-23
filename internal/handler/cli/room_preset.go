package cli

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/application"
	"github.com/jonatak/go-bailup/internal/domain"
)

type RoomPreset struct {
	Preset domain.ThermostatPreset `arg:"" enum:"eco,comfort" help:"eco or comfort mode"`
	RoomTarget
}

func (r *RoomPreset) Run(service *application.HVACService) error {

	_, err := service.SetRoomPreset(r.Name, r.Preset)
	if err != nil {
		return err
	}
	fmt.Printf("%s mode is now: %s\n", r.Name, r.Preset)

	return nil
}
