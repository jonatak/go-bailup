package command

import "encoding/json"

type JSONCommand interface {
	ToJSON() ([]byte, error)
}

type EmptyCommand struct{}

func (*EmptyCommand) ToJSON() ([]byte, error) {
	payload := map[string]any{}
	return json.Marshal(payload)
}
