package cli

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/app"
)

type Status struct{}

func (s *Status) Run(appCtx *app.AppContext) error {
	state, err := appCtx.HVACService.CurrentState()

	if err != nil {
		return err
	}

	fmt.Println(formatHVACSystem(state))
	return nil
}
