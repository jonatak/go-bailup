package cli

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/application"
)

type Status struct{}

func (s *Status) Run(service *application.HVACService) error {
	state, err := service.CurrentState()

	if err != nil {
		return err
	}

	fmt.Println(formatHVACSystem(state))
	return nil
}
