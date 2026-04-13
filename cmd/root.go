package cmd

import kongcompletion "github.com/jotaen/kong-completion"

type CLI struct {
	Status     Status                    `cmd:"" help:"Show Thermostats status"`
	HvacMode   HvacMode                  `cmd:"" help:"Set hvac mode"`
	Completion kongcompletion.Completion `cmd:"" help:"Outputs shell code for initialising tab completions"`
}
