# Go-Bailup

Go-Bailup is a small Go application for controlling Bailup / Baillconnect thermostats from the terminal and from Home Assistant over MQTT.

The project logs into the Baillconnect web interface, keeps the authenticated session in a cookie jar, reads the current regulation state, and sends command payloads for HVAC mode, room power, presets, and temperature setpoints.

## Status

This is a personal project, but the main flows are now usable:

- Fetch and display current thermostat state.
- Change global HVAC mode.
- List available rooms.
- Turn a room thermostat on or off.
- Switch a room between `eco` and `comfort`.
- Set, increase, or decrease room temperature setpoints.
- Run an MQTT/Home Assistant bridge with discovery, command handling, and state publishing.
- Generate shell completion scripts.

## Install

### GitHub Releases (latest binary)

On Linux or macOS (amd64 or arm64), install the latest published release into `~/.local/bin` (or `XDG_BIN_HOME` if set):

```sh
curl -fsSL https://raw.githubusercontent.com/jonatak/go-bailup/main/scripts/install-latest-release.sh | bash
```

Add the install directory to your `PATH` if it is not already, for example:

```sh
export PATH="$HOME/.local/bin:$PATH"
```

### Build from source

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

Install into your Go binary directory:

```sh
make install
```

## Configuration

The application reads Bailup credentials and regulation id from environment variables:

```sh
export BAILUP_EMAIL="you@example.com"
export BAILUP_PASS="your-password"
export BAILUP_REGULATION="your-regulation-id"
```

For MQTT / Home Assistant mode, also set:

```sh
export MQTT_HOST="mqtt.example.local"
export MQTT_PORT="1883"
export MQTT_USERNAME="mqtt-user"
export MQTT_PASSWORD="mqtt-password"
export MQTT_TOPIC_PREFIX="custom_bailup"
export MQTT_CLIENT_ID="go-bailup"
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

Run the MQTT / Home Assistant bridge:

```sh
bailup serve
```

The bridge:

- subscribes to command topics under `MQTT_TOPIC_PREFIX`
- publishes Home Assistant MQTT discovery payloads
- publishes thermostat and general state
- retries MQTT and Bailup connections when they drop

## Completion

Generate shell completion code with:

```sh
bailup completion
```

Use the generated output according to your shell setup.

## Architecture

The project is split into a few focused packages:

- `cmd/bailup`: installable CLI entrypoint.
- `internal/bootstrap`: application initialization and environment loading.
- `internal/domain`: HVAC aggregate, thermostat behavior, and setpoint rules.
- `internal/application`: use-case orchestration, inbound intents, target resolution, and outbound gateway port.
- `internal/infrastructure/cli`: Kong-based CLI commands and terminal formatting.
- `internal/infrastructure/mqtt`: MQTT command/message handling.
- `internal/infrastructure/bailup`: authenticated Baillconnect gateway, login flow, state mapping, and resolved-intent to command mapping.
- `internal/infrastructure/bailup/command`: JSON payload types sent to Baillconnect.
- `internal/infrastructure/bailup/model`: Baillconnect API DTOs and mode conversions.

The main flows are:

```text
CLI / MQTT message
  -> application.Intent
  -> application.HVACService.ApplyIntent
  -> domain.HVACSystem
  -> application.ResolvedIntent
  -> application.HVACSystemGateway
  -> infrastructure/bailup.Gateway
  -> Baillconnect HTTP API
```

Temperature requests use a two-step model in `application`:

- `Intent`: inbound request semantics such as `current` mode/preset or delta-based temperature changes.
- `ResolvedIntent`: gateway-ready operations with fully resolved domain values.

That keeps request interpretation in the application layer and keeps the Bailup adapter focused on vendor-specific command mapping.

The MQTT runtime is split into:

- `Handler`: MQTT transport, topic registration, discovery publish, and state publish.
- `Processor`: reconnect loop, worker loop, periodic refresh, and application service calls.

Home Assistant discovery uses:

- `homeassistant/climate/general/config`
- `homeassistant/climate/th_<id>/config`

Command and state topics use the configured `MQTT_TOPIC_PREFIX` and stable thermostat IDs such as `th_9152`.

## Development

Build everything:

```sh
go build ./...
```

Run tests and static checks:

```sh
go test ./...
go vet ./...
```

Run the CLI locally:

```sh
make build
./bin/bailup status
```

Run the MQTT bridge locally:

```sh
make build
./bin/bailup serve
```

## Deployment

For this project, the simplest deployment is usually a single binary under `systemd`.

Example unit:

```ini
[Unit]
Description=Go Bailup MQTT bridge
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=go-bailup
Group=go-bailup
EnvironmentFile=/etc/default/go-bailup
ExecStart=/usr/local/bin/bailup serve
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Example environment file:

```sh
BAILUP_EMAIL=you@example.com
BAILUP_PASS=your-password
BAILUP_REGULATION=your-regulation-id
MQTT_HOST=mqtt.example.local
MQTT_PORT=1883
MQTT_USERNAME=mqtt-user
MQTT_PASSWORD=mqtt-password
MQTT_TOPIC_PREFIX=custom_bailup
MQTT_CLIENT_ID=go-bailup
```

This keeps deployment simple and makes logs available through `journalctl`.


## Libraries

- [Kong](https://github.com/alecthomas/kong) for CLI parsing.
- [kong-completion](https://github.com/jotaen/kong-completion) for shell completion.
- [htmlquery](https://github.com/antchfx/htmlquery) for extracting login tokens.
