package cli

import (
	"context"
	"fmt"

	"github.com/jonatak/go-bailup/internal/application"
	"github.com/jonatak/go-bailup/internal/domain"
)

type RoomPreset struct {
	Preset domain.ThermostatPreset `arg:"" enum:"eco,comfort" help:"eco or comfort mode"`
	RoomTarget
}

func (r *RoomPreset) Run(ctx context.Context, service *application.HVACService) error {
	_, err := service.ApplyIntent(ctx, application.SetRoomPresetIntent{
		Room:   r.Name,
		Preset: r.Preset,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s mode is now: %s\n", r.Name, r.Preset)

	return nil
}
