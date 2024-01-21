package infrastructure

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"release-notifier/api"
)

var clientInstance *redis.Client

func GetClient() *redis.Client {
	if clientInstance == nil {
		clientInstance = setupClient()
	}

	return clientInstance
}

func setupClient() *redis.Client {
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

func GetCacheKey(release api.Release) string {
	return fmt.Sprintf("%s_%s", release.Url, release.TagName)
}
