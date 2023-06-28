package redis

import (
	"crypto/tls"

	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(address string, username string, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:      address,
		Username:  username,
		Password:  password,
		DB:        0,
		TLSConfig: &tls.Config{},
	})
}
