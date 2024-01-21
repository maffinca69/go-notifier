package cache

import (
	"context"
	"fmt"
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
	var key = infrastructure.GetCacheKey(release)

	client := infrastructure.GetClient()
	if err := client.Set(ctx, key, release.TagName, 0).Err(); err != nil {
		panic("Error save version to cache")
	}
}

func IsExists(release api.Release) bool {
	ctx := context.Background()

	client := infrastructure.GetClient()
	key := getCacheKey(release)
	if exists, _ := client.Exists(ctx, key).Result(); exists == 0 {
		client.Set(ctx, key, release.TagName, 0)
		return false
	}

	return true
}

func getCacheKey(release api.Release) string {
	return fmt.Sprintf("%s_%s", release.Url, release.TagName)
}
