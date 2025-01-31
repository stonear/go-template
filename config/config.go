package config

import (
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type Config struct {
	AppName    string `mapstructure:"APP_NAME"`
	AppVersion string `mapstructure:"APP_VERSION"`
	AppEnv     string `mapstructure:"APP_ENV"`
	AppKey     string `mapstructure:"APP_KEY"`
	AppDebug   bool   `mapstructure:"APP_DEBUG"`
	AppUrl     string `mapstructure:"APP_URL"`

	DbConnection string `mapstructure:"DB_CONNECTION"`
	DbHost       string `mapstructure:"DB_HOST"`
	DbPort       int    `mapstructure:"DB_PORT"`
	DbDatabase   string `mapstructure:"DB_DATABASE"`
	DbUsername   string `mapstructure:"DB_USERNAME"`
	DbPassword   string `mapstructure:"DB_PASSWORD"`
	DbTz         string `mapstructure:"DB_TZ"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     int    `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDb       int    `mapstructure:"REDIS_DB"`

	EnableTelemetry          bool   `mapstructure:"ENABLE_TELEMETRY"`
	OtelExporterOtlpEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OtelExporterOtlpHeaders  string `mapstructure:"OTEL_EXPORTER_OTLP_HEADERS"`
}

func New(log *otelzap.Logger) *Config {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		log.Warn("[Config] failed to load configuration file, use environment variables instead", zap.Error(err))
	} else {
		log.Info("[Config] configuration file loaded", zap.String("config file", viper.ConfigFileUsed()))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("[Config] error decoding configuration file", zap.Error(err))
	}

	return &config
}
