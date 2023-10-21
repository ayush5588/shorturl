package db

import (
	"errors"

	"github.com/go-redis/redis"
)

var (
	// ErrRedisNotConnected ...
	ErrRedisNotConnected = errors.New("cannot establish connection to redis database")
)

// NewRedisConnection ...
func NewRedisConnection() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, ErrRedisNotConnected
	}

	return client, nil

}
