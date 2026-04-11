package command

type Kind string

const (
	PresetMode  Kind = "preset_mode"
	Mode        Kind = "mode"
	Temperature Kind = "temperature"
)

type JSONCommand interface {
	ToJSON() ([]byte, error)
}
