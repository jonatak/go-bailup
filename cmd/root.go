package cmd

var CLI struct {
	Status      Status   `cmd:"" help:"Show Thermostats status"`
	SetHvacMode HvacMode `cmd:"" help:"Set hvac mode"`
}
