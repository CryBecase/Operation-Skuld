package xredis

import (
	"context"
	"runtime"

	"github.com/go-redis/redis/v8"
)

func New(c Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.DB,
		MinIdleConns: 10 * runtime.NumCPU(),
		PoolSize:     50 * runtime.NumCPU(),
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return client
}
