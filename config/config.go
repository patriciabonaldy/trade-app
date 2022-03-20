package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/patriciabonaldy/zero/internal/platform/logger"

	"github.com/spf13/viper"
)

// Config is the struct we use to configure the binaries
type Config struct {
	// Clients
	TradingPairs []string `mapstructure:"TRADING_PAIR"`
	ExchangeURL  string   `mapstructure:"EXCHANGE_URL"`
	MaxSize      int      `mapstructure:"MAX_SIZE"`
	Log          logger.Logger
}

// NewConfig returns a new configuration, and attempts to load
// config from file system.
func NewConfig() (*Config, error) {
	env := strings.ToLower(os.Getenv("GO_ENVIRONMENT"))
	viper.SetConfigFile(fmt.Sprintf("%s.env", env))
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var vipErr viper.ConfigFileNotFoundError
		if ok := errors.As(err, &vipErr); ok {
			log.Fatalln(fmt.Errorf("config file not found. %w", err))
		} else {
			log.Fatalln(fmt.Errorf("unexpected error loading config file. %w", err))
		}
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalln(fmt.Errorf("failed to unmarshal config. %w", err))
	}

	config.Log = logger.New()

	return config, nil
}
