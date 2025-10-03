package config

import (
	"errors"
)

const (
	DB_PORT int = 5432
)

type DatabaseConfiguration struct {
	Host     string `yaml:"host,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	Database string `yaml:"database,omitempty"`
	Port     int    `yaml:"port,omitempty"`
}

func (c *DatabaseConfiguration) SetDefaults() {
	if c.Port == 0 {
		c.Port = DB_PORT
	}
}

func (c *DatabaseConfiguration) Validate() error {
	if len(c.User) == 0 {
		return errors.New("db user is required")
	}

	if len(c.Password) == 0 {
		return errors.New("db password is required")
	}

	if len(c.Database) == 0 {
		return errors.New("database name is required")
	}

	if c.Port == 0 {
		return errors.New("db port is required")
	}

	return nil
}
