# Baillconnect to MQTT

Baillconnect to MQTT is a small Go application for exposing Bailup / Baillconnect thermostats to Home Assistant over MQTT.

The project logs into the Baillconnect web interface, keeps the authenticated session in a cookie jar, reads the current regulation state, and sends command payloads for HVAC mode, room power, presets, and temperature setpoints.

## Status

This is a personal project, but the main flows are now usable:

- Run an MQTT/Home Assistant bridge with discovery, command handling, and state publishing.
- Change global HVAC mode from Home Assistant.
- Turn a room thermostat on or off from Home Assistant.
- Switch a room between `eco` and `comfort` from Home Assistant.
- Set room temperature setpoints from Home Assistant.

## Install

### Home Assistant Add-on

This repository can be added directly as a Home Assistant add-on repository:

```text
https://github.com/jonatak/baillconnect-to-mqtt
```

In Home Assistant OS or Supervised, open **Settings > Add-ons > Add-on Store**, add this repository URL, then install **Baillconnect to MQTT**.

Add-on documentation lives in [`baillconnect-to-mqtt/DOCS.md`](baillconnect-to-mqtt/DOCS.md).

### GitHub Releases (latest binary)

On Linux or macOS (amd64, arm64, or armv7), install the latest published release into `~/.local/bin` (or `XDG_BIN_HOME` if set):

```sh
curl -fsSL https://raw.githubusercontent.com/jonatak/baillconnect-to-mqtt/main/scripts/install-latest-release.sh | bash
```

Add the install directory to your `PATH` if it is not already, for example:

```sh
export PATH="$HOME/.local/bin:$PATH"
```

### Build from source

Clone and build:

```sh
git clone https://github.com/jonatak/baillconnect-to-mqtt.git
cd baillconnect-to-mqtt
make build
```

The binary is expected at:

```sh
./bin/baillconnect-to-mqtt
```

Install into your Go binary directory:

```sh
make install
```

## Configuration

The Home Assistant add-on passes configuration through `/data/options.json`.

For standalone binary usage, the same JSON shape can be provided with `--config`:

```json
{
  "baillconnect": {
    "email": "you@example.com",
    "password": "your-password",
    "regulation": "your-regulation-id"
  },
  "mqtt": {
    "host": "mqtt.example.local",
    "port": 1883,
    "username": "mqtt-user",
    "password": "mqtt-password",
    "topic_prefix": "custom_bailup",
    "client_id": "baillconnect-to-mqtt"
  },
  "poll_interval_seconds": 30
}
```

Environment variables are also supported:

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
export MQTT_CLIENT_ID="baillconnect-to-mqtt"
```

If you use `direnv`, put them in `.envrc` locally. Do not commit real credentials.

## Usage

Run the MQTT / Home Assistant bridge:

```sh
baillconnect-to-mqtt
```

Run with a config file:

```sh
baillconnect-to-mqtt --config ./options.json
```

The bridge:

- subscribes to command topics under `MQTT_TOPIC_PREFIX`
- publishes Home Assistant MQTT discovery payloads
- publishes thermostat and general state
- retries MQTT and Bailup connections when they drop

## Architecture

The project is split into a few focused packages:

- `cmd/baillconnect-to-mqtt`: installable application entrypoint.
- `internal/bootstrap`: application initialization and environment loading.
- `internal/domain`: HVAC aggregate, thermostat behavior, and setpoint rules.
- `internal/application`: use-case orchestration, inbound intents, target resolution, and outbound gateway port.
- `internal/mqtt`: MQTT command/message handling.
- `internal/bailup`: authenticated Baillconnect gateway, login flow, state mapping, and resolved-intent to command mapping.
- `internal/bailup/command`: JSON payload types sent to Baillconnect.
- `internal/bailup/model`: Baillconnect API DTOs and mode conversions.

The main flows are:

```text
MQTT message
  -> application.Intent
  -> application.HVACService.ApplyIntent
  -> domain.HVACSystem
  -> application.ResolvedIntent
  -> application.HVACSystemGateway
  -> bailup.Gateway
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

Run the bridge locally:

```sh
make build
./bin/baillconnect-to-mqtt
```

Run with a local config file:

```sh
./bin/baillconnect-to-mqtt --config ./options.json
```

## Deployment

For this project, the simplest deployment is usually a single binary under `systemd`.

Example unit:

```ini
[Unit]
Description=Baillconnect MQTT bridge
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=baillconnect-to-mqtt
Group=baillconnect-to-mqtt
EnvironmentFile=/etc/default/baillconnect-to-mqtt
ExecStart=/usr/local/bin/baillconnect-to-mqtt
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
MQTT_CLIENT_ID=baillconnect-to-mqtt
```

This keeps deployment simple and makes logs available through `journalctl`.


## Libraries

- [htmlquery](https://github.com/antchfx/htmlquery) for extracting login tokens.
- [viper](https://github.com/spf13/viper) for configuration management.
- [Paho Mqtt Golang](https://github.com/eclipse/paho.mqtt.golang) to manage mqtt connection.
