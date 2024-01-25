package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/xlab/closer"
	"release-notifier/cache"
	"release-notifier/github"
	"release-notifier/infrastructure"
	"release-notifier/telegram"
	"sync"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func main() {
	c := cron.New()

	if _, err := c.AddFunc(infrastructure.GetConfig().CronExpression, func() { checkUpdates() }); err != nil {
		panic("Error start schedule function")
	}

	c.Start()

	closer.Hold()
}

func checkUpdates() {
	var wg sync.WaitGroup

	for _, repo := range infrastructure.GetConfig().Repository {
		wg.Add(1)
		repo := repo
		go func() {
			defer wg.Done()
			notifyIfNeeded(repo)
		}()
	}

	wg.Wait()
}

func notifyIfNeeded(repo infrastructure.RepositoryConfig) {
	fmt.Println("Checking new version for:", repo.Name)
	var latestRelease = github.GetLatestRelease(repo.Url)

	if latestRelease == nil {
		fmt.Println(fmt.Sprintf("%s: not found latest release. Skip", repo.Name))
		return
	}

	if isAvailableNewVersion(repo, latestRelease) {
		telegram.Notify(latestRelease, repo.Name)
		cache.Save(repo.Name, latestRelease.TagName)
	}
}

func isAvailableNewVersion(repo infrastructure.RepositoryConfig, release *github.Release) bool {
	if cache.IsExists(repo.Name) == false {
		cache.Save(repo.Name, release.TagName)
		return false
	}

	version := cache.GetCurrentVersion(repo.Name)
	var isAvailable = version != release.TagName

	if isAvailable {
		fmt.Println(fmt.Sprintf("%s: new version available! %s", repo.Name, version))
	} else {
		fmt.Println(fmt.Sprintf("%s: New version not available. Current version %s", repo.Name, version))
	}

	return isAvailable
}
