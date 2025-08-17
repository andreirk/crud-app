package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB Postgres

	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

func (c *Config) mapPostgres() {
	c.DB.Host = c.mapField("host")
	c.DB.Port = c.mapField("port")
	c.DB.User = c.mapField("user")
	c.DB.Pass = c.mapField("pass")
	c.DB.Name = c.mapField("name")
}

func (c *Config) mapField(field string) string {
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

	cfg.mapPostgres()

	return cfg, nil
}
