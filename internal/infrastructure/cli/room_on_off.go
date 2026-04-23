package cli

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/application"
)

type RoomOn struct {
	RoomTarget
}

type RoomOff struct {
	RoomTarget
}

func (o *RoomOn) Run(service *application.HVACService) error {
	_, err := service.TurnRoomOn(o.Name)
	if err != nil {
		return err
	}

	fmt.Printf("%s is now: on\n", o.Name)
	return nil
}

func (o *RoomOff) Run(service *application.HVACService) error {
	_, err := service.TurnRoomOff(o.Name)
	if err != nil {
		return err
	}

	fmt.Printf("%s is now: off\n", o.Name)
	return nil
}
