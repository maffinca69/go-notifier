package cache

import (
	"context"
	"release-notifier/infrastructure"
)

func GetCurrentVersion(repository string) string {
	ctx := context.Background()

	client := infrastructure.RedisClient()
	version, _ := client.Get(ctx, repository).Result()

	return version
}

func Save(repository string, version string) {
	ctx := context.Background()

	client := infrastructure.RedisClient()
	if err := client.Set(ctx, repository, version, 0).Err(); err != nil {
		panic("Error save version to cache")
	}
}

func IsExists(repository string) bool {
	ctx := context.Background()

	client := infrastructure.RedisClient()
	exists, err := client.Exists(ctx, repository).Result()
	if err != nil {
		panic(err)
	}

	return exists == 1
}
