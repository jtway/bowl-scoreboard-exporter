package bowl_exporter

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Endpoint      string        `mapstructure:"api_endpoint"`
	FetchInterval time.Duration `mapstructure:"fetch_interval"`
	Prom          struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"prom"`
}

func ReadInConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/bowl-scoreboard-exporter/")
	viper.AddConfigPath("$HOME/.bowl-scoreboard-exporter")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into config struct, %w", err)
	}

	if len(config.Endpoint) == 0 {
		config.Endpoint = "http://site.api.espn.com/apis/site/v2/sports/football/college-football/scoreboard"
	}
	if config.FetchInterval == time.Duration(0) {
		config.FetchInterval = 60 * time.Second
	}
	if config.Prom.Port == 0 {
		config.Prom.Port = 9180
	}
	return config, nil
}
