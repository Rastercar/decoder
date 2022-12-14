package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	Debug             bool `yaml:"debug" env:"APP_DEBUG"`
	MaxInvalidPackets int  `yaml:"max_invalid_packets" env:"MAX_INVALID_PACKETS"`
}

type TracerConfig struct {
	Url         string `env-required:"true" yaml:"url" env:"TRACER_URL"`
	ServiceName string `env-required:"true" yaml:"service_name" env:"TRACER_SERVICE_NAME"`
}

type RmqConfig struct {
	Url               string `env-required:"true" yaml:"url" env:"RMQ_URL"`
	Exchange          string `env-required:"true" yaml:"exchange" env:"RMQ_EXCHANGE"`
	ReconnectWaitTime int    `env-required:"true" yaml:"reconnect_wait_time" env:"RMQ_RECONNECT_WAIT_TIME"`
}

type Config struct {
	App    AppConfig    `yaml:"app"`
	Rmq    RmqConfig    `yaml:"rmq"`
	Tracer TracerConfig `yaml:"tracer"`
}

func Parse() (*Config, error) {
	var cfgFilePath = flag.String("config-file", "/etc/config.yml", "A filepath to the yml file containing the microservice configuration")
	flag.Parse()

	cfg := &Config{}

	if err := cleanenv.ReadConfig(*cfgFilePath, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
