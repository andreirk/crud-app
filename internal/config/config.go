package config

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB     Postgres
	Salt   string
	Secret []byte

	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	Auth struct {
		TokenTTL time.Duration `mapstructure:"token_ttl"`
	} `mapstructure:"auth"`
}

type Postgres struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func (c *Config) mapPostgres() error {
	c.DB.Host = c.getField("db_host")
	c.DB.Port = c.getField("db_port")
	c.DB.User = c.getField("db_user")
	c.DB.Pass = c.getField("db_pass")
	c.DB.Name = c.getField("db_name")

	if c.DB.Host == "" || c.DB.Port == "" ||
		c.DB.User == "" || c.DB.Pass == "" || c.DB.Name == "" {
		return errors.New("unable to map all fields")
	}

	return nil
}

func (c *Config) getField(field string) string {
	if val, ok := viper.Get(field).(string); ok {
		return val
	}

	return ""
}

func (c *Config) getSalt() error {
	val, ok := viper.Get("hash_salt").(string)
	if ok {
		c.Salt = val
		return nil
	}

	return errors.New("unable to retrieve salt")
}

func (c *Config) getSecret() error {
	val, ok := viper.Get("jwt_secret").(string)
	if ok {
		c.Secret = []byte(val)
		return nil
	}

	return errors.New("unable to retrieve secret")
}

func New(dir, file string) (*Config, error) {
	cfg := new(Config)

	viper.AddConfigPath(dir)
	viper.SetConfigName(file)
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

	if err := cfg.getSalt(); err != nil {
		return nil, err
	}

	if err := cfg.getSecret(); err != nil {
		return nil, err
	}

	return cfg, nil
}
