package cli

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/app"
	"github.com/jonatak/go-bailup/internal/domain"
)

type HVACMode struct {
	Mode domain.HVACSystemMode `arg:"" enum:"off,cool,heat,dry,fan-only"`
}

func (s *HVACMode) Run(appCtx *app.AppContext) error {

	state, err := appCtx.HVACService.SetMode(s.Mode)
	if err != nil {
		return err
	}

	fmt.Printf("HVAC mode is now: %s\n", state.Mode())

	return nil
}
