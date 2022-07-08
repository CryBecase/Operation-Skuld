package xredis

import (
	"time"
)

type Config struct {
	Addr         string        `mapstructure:"addr"`
	DB           int           `mapstructure:"db"`
	Password     string        `mapstructure:"password"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}
