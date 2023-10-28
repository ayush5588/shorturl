package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

var (
	// ErrRedisNotConnected ...
	ErrRedisNotConnected = errors.New("cannot establish connection to redis database")
)

var (
	rHost = os.Getenv("REDIS_HOST")
	rPort = os.Getenv("REDIS_PORT")
)

func getRedisAddr() string {
	// If rHost or rPort are not set like when directly running in local system and not through docker
	if rHost == "" {
		rHost = "localhost"
	}
	if rPort == "" {
		rPort = "6379"
	}

	redisAddr := fmt.Sprintf("%s:%s", rHost, rPort)

	return redisAddr
}

// NewRedisConnection ...
func NewRedisConnection() (*redis.Client, error) {
	redisAddr := getRedisAddr()

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, ErrRedisNotConnected
	}

	return client, nil

}
