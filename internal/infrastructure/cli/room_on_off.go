package cli

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/app"
)

type RoomOn struct {
	RoomTarget
}

type RoomOff struct {
	RoomTarget
}

func (o *RoomOn) Run(appCtx *app.AppContext) error {
	_, err := appCtx.HVACService.TurnRoomOn(o.Name)
	if err != nil {
		return err
	}

	fmt.Printf("%s is now: on\n", o.Name)
	return nil
}

func (o *RoomOff) Run(appCtx *app.AppContext) error {
	_, err := appCtx.HVACService.TurnRoomOff(o.Name)
	if err != nil {
		return err
	}

	fmt.Printf("%s is now: off\n", o.Name)
	return nil
}
