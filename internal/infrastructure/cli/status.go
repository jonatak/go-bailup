package cli

import (
	"context"
	"fmt"

	"github.com/jonatak/go-bailup/internal/application"
)

type Status struct{}

func (s *Status) Run(ctx context.Context, service *application.HVACService) error {
	state, err := service.CurrentState(ctx)

	if err != nil {
		return err
	}

	fmt.Println(formatHVACSystem(state))
	return nil
}
