# Go-Bailup

Go-Bailup is a small Go CLI for controlling Bailup / Baillconnect thermostats from the terminal.

The project logs into the Baillconnect web interface, keeps the authenticated session in a cookie jar, reads the current regulation state, and sends command payloads for HVAC mode, room power, presets, and temperature setpoints.

## Status

This is an early-stage personal project, but the main CLI flow is already usable:

- Fetch and display current thermostat state.
- Change global HVAC mode.
- List available rooms.
- Turn a room thermostat on or off.
- Switch a room between `eco` and `comfort`.
- Set, increase, or decrease room temperature setpoints.
- Generate shell completion scripts.

## Install

Clone and build:

```sh
git clone https://github.com/jonatak/go-bailup.git
cd go-bailup
make build
```

The binary is expected at:

```sh
./bin/bailup
```

## Configuration

The CLI reads credentials and regulation id from environment variables:

```sh
export BAILUP_EMAIL="you@example.com"
export BAILUP_PASS="your-password"
export BAILUP_REGULATION="your-regulation-id"
```

If you use `direnv`, put them in `.envrc` locally. Do not commit real credentials.

## Usage

Show full thermostat status:

```sh
bailup status
```

Set the global HVAC mode:

```sh
bailup hvac-mode heat
bailup hvac-mode cool
bailup hvac-mode off
```

List rooms:

```sh
bailup room list
```

Turn a room thermostat on or off:

```sh
bailup room on "Living Room"
bailup room off "Bedroom"
```

Switch a room preset:

```sh
bailup room preset comfort "Living Room"
bailup room preset eco "Bedroom"
```

Set a room temperature:

```sh
bailup room temp set "Living Room" 20
```

Increase or decrease a room temperature:

```sh
bailup room temp up "Living Room"
bailup room temp down "Living Room"
bailup room temp up "Living Room" --by 0.5
```

Target a specific preset or HVAC setpoint:

```sh
bailup room temp set "Living Room" 19 --preset eco
bailup room temp set "Living Room" 21 --preset comfort --mode heat
bailup room temp up "Living Room" --by 1 --preset eco --mode cool
```

`--preset current` and `--mode current` are the defaults. Use explicit values when the current HVAC mode is not enough to identify the setpoint you want to modify.

## Completion

Generate shell completion code with:

```sh
bailup completion
```

Use the generated output according to your shell setup.

## Architecture

The project is split into a few focused packages:

- `cmd`: Kong-based CLI commands.
- `internal/app`: application initialization and environment loading.
- `internal/bailup`: authenticated Baillconnect client, login flow, request execution, and command builders.
- `internal/bailup/command`: JSON payload types sent to Baillconnect.
- `internal/bailup/model`: API models, modes, and human-readable formatting.

## Development

Build everything:

```sh
go build ./...
```

Run the CLI locally:

```sh
make build
./bin/bailup status
```

## Roadmap

- [ ] Add unit tests for command payload generation.
- [ ] Add tests for mode parsing and temperature targeting.
- [ ] Add MQTT / Home Assistant integration.
- [ ] Reduce HTTP headers to the minimum required by Baillconnect.
- [ ] Add a real server-side session check if needed.

## Libraries

- [Kong](https://github.com/alecthomas/kong) for CLI parsing.
- [kong-completion](https://github.com/jotaen/kong-completion) for shell completion.
- [htmlquery](https://github.com/antchfx/htmlquery) for extracting login tokens.
