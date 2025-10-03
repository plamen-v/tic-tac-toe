package config

import "errors"

type ServerConfiguration struct {
	Port int `yaml:"port,omitempty"`
}

func (c *ServerConfiguration) SetDefaults() {
}

func (c *ServerConfiguration) Validate() error {
	if c.Port == 0 {
		return errors.New("application port is invalid")
	}

	return nil
}
