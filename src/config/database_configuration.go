package config

import "fmt"

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
		return fmt.Errorf("todo")
	}

	if len(c.Password) == 0 {
		return fmt.Errorf("todo")
	}

	if len(c.Database) == 0 {
		return fmt.Errorf("todo")
	}

	if c.Port == 0 {
		return fmt.Errorf("todo")
	}

	return nil
}
