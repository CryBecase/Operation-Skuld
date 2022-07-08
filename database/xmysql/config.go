package xmysql

import (
	"time"
)

type Config struct {
	Driver   string        `mapstructure:"driver"`
	Source   string        `mapstructure:"source"`
	Idle     int           `mapstructure:"idle"`
	Open     int           `mapstructure:"open"`
	IdleTime time.Duration `mapstructure:"idle_time"`
}
