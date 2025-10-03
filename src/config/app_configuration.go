package config

import (
	"errors"
	"flag"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type (
	AppMode  string
	LogLevel string
)

const (
	ProductionAppMode  AppMode  = "prod"
	DevelopmentAppMode AppMode  = "dev"
	DebugLogLevel      LogLevel = "debug"
	InfoLogLevel       LogLevel = "info"
)

var (
	configPath string
	config     AppConfiguration
	once       sync.Once
)

func init() {
	flag.StringVar(
		&configPath,
		"config",
		"",
		"Path to the configuration file")
}

type AppConfiguration struct {
	AppName  string                `yaml:"appName,omitempty"`
	AppMode  AppMode               `yaml:"appMode,omitempty"`
	LogLevel LogLevel              `yaml:"logLevel,omitempty"`
	Secret   string                `yaml:"secret,omitempty"`
	Server   ServerConfiguration   `yaml:"server"`
	Database DatabaseConfiguration `yaml:"database"`
}

func (c *AppConfiguration) SetDefaults() {
	if len(c.AppMode) == 0 {
		c.AppMode = ProductionAppMode
	}
	if len(c.LogLevel) == 0 {
		if c.AppMode == DevelopmentAppMode {
			c.LogLevel = DebugLogLevel
		} else {
			c.LogLevel = InfoLogLevel
		}
	}

	c.Server.SetDefaults()
	c.Database.SetDefaults()
}

func (c *AppConfiguration) Validate() error {
	if len(c.AppName) == 0 {
		return errors.New("application name is required")
	}

	if len(c.Secret) == 0 {
		return errors.New("application secret is required")
	}

	if err := c.Server.Validate(); err != nil {
		return err
	}

	if err := c.Database.Validate(); err != nil {
		return err
	}

	return nil
}

func GetConfig() (*AppConfiguration, error) {
	once.Do(func() {
		data, err := os.ReadFile(configPath)
		if err != nil {
			panic(err)
		}

		err = yaml.NewDecoder(strings.NewReader(os.ExpandEnv(string(data)))).Decode(&config)
		if err != nil {
			panic(err)
		}
	})

	config.SetDefaults()
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}
