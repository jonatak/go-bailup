package command

import "encoding/json"

type Kind string

const (
	PresetMode  Kind = "preset_mode"
	Mode        Kind = "mode"
	Temperature Kind = "temperature"
)

type JSONCommand interface {
	ToJSON() ([]byte, error)
}

type EmptyCommand struct{}

func (*EmptyCommand) ToJSON() ([]byte, error) {
	payload := map[string]any{}
	return json.Marshal(payload)
}
