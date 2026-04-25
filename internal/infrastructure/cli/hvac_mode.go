package cli

import (
	"context"
	"fmt"

	"github.com/jonatak/go-bailup/internal/application"
	"github.com/jonatak/go-bailup/internal/domain"
)

type HVACMode struct {
	Mode domain.HVACSystemMode `arg:"" enum:"off,cool,heat,dry,fan-only"`
}

func (s *HVACMode) Run(ctx context.Context, service *application.HVACService) error {
	state, err := service.ApplyIntent(ctx, application.SetModeIntent{
		Mode: s.Mode,
	})
	if err != nil {
		return err
	}

	fmt.Printf("HVAC mode is now: %s\n", state.Mode())

	return nil
}
