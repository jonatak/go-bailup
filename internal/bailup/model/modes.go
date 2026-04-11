package model

import (
	"fmt"
	"strings"
)

type UCMode int

const (
	UCModeOff UCMode = iota
	UCModeCool
	UCModeHeat
	UCModeDry
	UCModeFanOnly
	UCModeAuto
)

type ThMode int

const (
	ThModeComfort ThMode = 1 + iota
	ThModeEco
)

func ThModeFromString(mode string) (ThMode, error) {
	switch strings.ToLower(mode) {
	case "comfort":
		return ThModeComfort, nil
	case "eco":
		return ThModeEco, nil
	default:
		return 0, fmt.Errorf("unsupported thermostat mode %q", mode)
	}
}

func UCModeFromString(mode string) (UCMode, error) {
	switch strings.ToLower(mode) {
	case "off":
		return UCModeOff, nil
	case "cool":
		return UCModeCool, nil
	case "heat":
		return UCModeHeat, nil
	case "dry":
		return UCModeDry, nil
	case "fan-only":
		return UCModeFanOnly, nil
	case "auto":
		return UCModeAuto, nil
	default:
		return 0, fmt.Errorf("unsupported unit mode %q", mode)
	}
}

func (m UCMode) String() string {
	switch m {
	case UCModeOff:
		return "off"
	case UCModeCool:
		return "cool"
	case UCModeHeat:
		return "heat"
	case UCModeDry:
		return "dry"
	case UCModeFanOnly:
		return "fan-only"
	case UCModeAuto:
		return "auto"
	default:
		return fmt.Sprintf("unknown(%d)", int(m))
	}
}

func (m ThMode) String() string {
	switch m {
	case ThModeComfort:
		return "comfort"
	case ThModeEco:
		return "eco"
	default:
		return fmt.Sprintf("unknown(%d)", int(m))
	}
}
