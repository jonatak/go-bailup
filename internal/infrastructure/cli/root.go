package cli

import kongcompletion "github.com/jotaen/kong-completion"

type CLI struct {
	Status     Status                    `cmd:"" help:"Show Thermostats status"`
	HVACMode   HVACMode                  `cmd:"" help:"Set HVAC mode"`
	Room       Room                      `cmd:"" help:"Manage rooms"`
	Completion kongcompletion.Completion `cmd:"" help:"Outputs shell code for initialising tab completions"`
}
