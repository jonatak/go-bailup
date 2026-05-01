package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Baillconnect BaillconnectConfig `mapstructure:"baillconnect"`
	MQTT         MQTTConfig         `mapstructure:"mqtt"`
	PollInterval int                `mapstructure:"poll_interval_seconds"`
}

type BaillconnectConfig struct {
	Email      string `mapstructure:"email"`
	Password   string `mapstructure:"password"`
	Regulation string `mapstructure:"regulation"`
}

type MQTTConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	TopicPrefix string `mapstructure:"topic_prefix"`
	ClientID    string `mapstructure:"client_id"`
}

func Load(configPath string) (Config, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)
	bindEnv(v)

	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("mqtt.host", "core-mosquitto")
	v.SetDefault("mqtt.port", 1883)
	v.SetDefault("mqtt.topic_prefix", "baillconnect")
	v.SetDefault("mqtt.client_id", "baillconnect-to-mqtt")
	v.SetDefault("poll_interval_seconds", 30)
}

func bindEnv(v *viper.Viper) {
	mustBindEnv(v, "baillconnect.email", "BAILUP_EMAIL", "BAILLCONNECT_EMAIL")
	mustBindEnv(v, "baillconnect.password", "BAILUP_PASS", "BAILUP_PASSWORD", "BAILLCONNECT_PASSWORD")
	mustBindEnv(v, "baillconnect.regulation", "BAILUP_REGULATION", "BAILLCONNECT_REGULATION")
	mustBindEnv(v, "mqtt.host", "MQTT_HOST")
	mustBindEnv(v, "mqtt.port", "MQTT_PORT")
	mustBindEnv(v, "mqtt.username", "MQTT_USERNAME")
	mustBindEnv(v, "mqtt.password", "MQTT_PASSWORD")
	mustBindEnv(v, "mqtt.topic_prefix", "MQTT_TOPIC_PREFIX")
	mustBindEnv(v, "mqtt.client_id", "MQTT_CLIENT_ID")
	mustBindEnv(v, "poll_interval_seconds", "POLL_INTERVAL_SECONDS")
}

func mustBindEnv(v *viper.Viper, key string, envVars ...string) {
	args := append([]string{key}, envVars...)
	if err := v.BindEnv(args...); err != nil {
		panic(err)
	}
}
