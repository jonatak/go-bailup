package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadReadsOptionsJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "options.json")
	err := os.WriteFile(path, []byte(`{
		"baillconnect": {
			"email": "user@example.com",
			"password": "secret",
			"regulation": "123"
		},
		"mqtt": {
			"host": "mqtt.local",
			"port": 1884,
			"username": "mqtt-user",
			"password": "mqtt-secret",
			"topic_prefix": "custom",
			"client_id": "client"
		},
		"poll_interval_seconds": 45
	}`), 0o600)
	require.NoError(t, err)

	cfg, err := Load(path)

	require.NoError(t, err)
	assert.Equal(t, "user@example.com", cfg.Baillconnect.Email)
	assert.Equal(t, "secret", cfg.Baillconnect.Password)
	assert.Equal(t, "123", cfg.Baillconnect.Regulation)
	assert.Equal(t, "mqtt.local", cfg.MQTT.Host)
	assert.Equal(t, 1884, cfg.MQTT.Port)
	assert.Equal(t, "mqtt-user", cfg.MQTT.Username)
	assert.Equal(t, "mqtt-secret", cfg.MQTT.Password)
	assert.Equal(t, "custom", cfg.MQTT.TopicPrefix)
	assert.Equal(t, "client", cfg.MQTT.ClientID)
	assert.Equal(t, 45, cfg.PollInterval)
}

func TestLoadFallsBackToLegacyEnvironment(t *testing.T) {
	t.Setenv("BAILUP_EMAIL", "env@example.com")
	t.Setenv("BAILUP_PASS", "env-secret")
	t.Setenv("BAILUP_REGULATION", "env-regulation")
	t.Setenv("MQTT_HOST", "env-mqtt.local")
	t.Setenv("MQTT_PORT", "1885")
	t.Setenv("MQTT_USERNAME", "env-mqtt-user")
	t.Setenv("MQTT_PASSWORD", "env-mqtt-secret")
	t.Setenv("MQTT_TOPIC_PREFIX", "env-prefix")
	t.Setenv("MQTT_CLIENT_ID", "env-client")
	t.Setenv("POLL_INTERVAL_SECONDS", "90")

	cfg, err := Load("")

	require.NoError(t, err)
	assert.Equal(t, "env@example.com", cfg.Baillconnect.Email)
	assert.Equal(t, "env-secret", cfg.Baillconnect.Password)
	assert.Equal(t, "env-regulation", cfg.Baillconnect.Regulation)
	assert.Equal(t, "env-mqtt.local", cfg.MQTT.Host)
	assert.Equal(t, 1885, cfg.MQTT.Port)
	assert.Equal(t, "env-mqtt-user", cfg.MQTT.Username)
	assert.Equal(t, "env-mqtt-secret", cfg.MQTT.Password)
	assert.Equal(t, "env-prefix", cfg.MQTT.TopicPrefix)
	assert.Equal(t, "env-client", cfg.MQTT.ClientID)
	assert.Equal(t, 90, cfg.PollInterval)
}
