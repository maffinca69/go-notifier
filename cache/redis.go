package cache

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

func GetClient() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	return client
}
