package cache

import (
	"context"
	"release-notifier/api"
	"release-notifier/infrastructure"
)

func GetCurrentVersion(release api.Release) string {
	ctx := context.Background()

	client := infrastructure.GetClient()
	key := getCacheKey(release)
	version, _ := client.Get(ctx, key).Result()

	return version
}

func Save(release api.Release) {
	ctx := context.Background()
	var key = getCacheKey(release)

	client := infrastructure.GetClient()
	if err := client.Set(ctx, key, release.TagName, 0).Err(); err != nil {
		panic("Error save version to cache")
	}
}

func IsExists(release api.Release) bool {
	ctx := context.Background()

	client := infrastructure.GetClient()
	key := getCacheKey(release)
	exists, err := client.Exists(ctx, key).Result()
	if err != nil {
		panic(err)
	}

	return exists == 1
}

func getCacheKey(release api.Release) string {
	return release.Url
}
