package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	DB Postgres

	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

func (c *Config) mapPostgres() error {
	c.DB.Host = c.getField("host")
	c.DB.Port = c.getField("port")
	c.DB.User = c.getField("user")
	c.DB.Pass = c.getField("pass")
	c.DB.Name = c.getField("name")

	if c.DB.Host == "" || c.DB.Port == "" ||
		c.DB.User == "" || c.DB.Pass == "" || c.DB.Name == "" {
		return errors.New("unable to retrieve all fields")
	}

	return nil
}

func (c *Config) getField(field string) string {
	if val, ok := viper.Get(field).(string); ok {
		return val
	}

	return ""
}

type Postgres struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func New(dir, file string) (*Config, error) {
	cfg := new(Config)

	viper.AddConfigPath(dir)
	viper.SetConfigName(file)
	viper.SetEnvPrefix("db")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if err := cfg.mapPostgres(); err != nil {
		return nil, err
	}

	return cfg, nil
}
