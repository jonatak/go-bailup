package command

import (
	"encoding/json"
)

type ModeCommand struct {
	Value int
}

func (m ModeCommand) ToJSON() ([]byte, error) {
	payload := map[string]any{}

	payload["uc_mode"] = m.Value
	return json.Marshal(payload)
}
