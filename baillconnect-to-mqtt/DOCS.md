# Baillconnect to MQTT

Bridge your Baillconnect HVAC system to Home Assistant using MQTT.

## Configuration

```yaml
baillconnect:
  email: ""
  password: ""
  regulation: ""

mqtt:
  host: "core-mosquitto"
  port: 1883
  username: ""
  password: ""
  topic_prefix: "baillconnect"
  client_id: "baillconnect-to-mqtt"

poll_interval_seconds: 30
```

## How It Works

The add-on runs the `baillconnect-to-mqtt` Go binary.

The binary is expected to:

- Read `/data/options.json`
- Connect to the Baillconnect API
- Connect to MQTT
- Publish Home Assistant MQTT discovery topics
- Publish HVAC state updates
- Listen for MQTT command topics

## Requirements

- Home Assistant OS or Supervised
- MQTT broker, with the Mosquitto add-on recommended
