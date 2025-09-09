package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	configPath string
	config     AppConfiguration
	once       sync.Once
)

const (
	APP_NAME string = "tic-tac-toe"
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
	Secret   string                `yaml:"secret,omitempty"`
	Server   ServerConfiguration   `yaml:"server"`
	Database DatabaseConfiguration `yaml:"database"`
}

func (c *AppConfiguration) SetDefaults() {
	if len(c.AppName) == 0 {
		c.AppName = APP_NAME
	}

	c.Server.SetDefaults()
	c.Database.SetDefaults()
}

// TODO!
func (c *AppConfiguration) Validate() error {
	if len(c.Secret) == 0 {
		return fmt.Errorf("todo")
	}

	if err := c.Server.Validate(); err != nil {
		panic(err)
	}

	if err := c.Database.Validate(); err != nil {
		panic(err)
	}

	return nil
}

func GetConfig() *AppConfiguration {
	once.Do(func() { //TODO!
		data, err := os.ReadFile(configPath)
		if err != nil {
			panic(err)
		}

		err = yaml.NewDecoder(strings.NewReader(os.ExpandEnv(string(data)))).Decode(&config) // TODO! detail
		if err != nil {
			panic(err)
		}
	})

	config.SetDefaults() //TODO!
	if err := config.Validate(); err != nil {
		panic(err)
	}

	return &config
}
