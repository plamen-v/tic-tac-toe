package config

import "fmt"

const (
	SERVER_PORT int    = 5005
	SERVER_HOST string = "localhost"
)

type ServerConfiguration struct {
	Host string `yaml:"host,omitempty"`
	Port int    `yaml:"port,omitempty"`
}

func (c *ServerConfiguration) SetDefaults() {
	if len(c.Host) == 0 {
		c.Host = SERVER_HOST
	}
	if c.Port == 0 {
		c.Port = SERVER_PORT
	}
}

func (c *ServerConfiguration) Validate() error {
	if len(c.Host) == 0 {
		return fmt.Errorf("todo")
	}

	if c.Port == 0 {
		return fmt.Errorf("todo")
	}

	return nil
}
