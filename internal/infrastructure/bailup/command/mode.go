package command

import (
	"encoding/json"

	"github.com/jonatak/go-bailup/internal/infrastructure/bailup/model"
)

type ModeCommand struct {
	Value int
}

func NewHVACModeCommand(mode string) (JSONCommand, error) {
	ucMode, err := model.UCModeFromString(mode)
	if err != nil {
		return nil, err
	}

	return ModeCommand{
		Value: int(ucMode),
	}, nil
}

func (m ModeCommand) ToJSON() ([]byte, error) {
	payload := map[string]any{}

	payload["uc_mode"] = m.Value
	return json.Marshal(payload)
}
