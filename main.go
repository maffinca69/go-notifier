package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/xlab/closer"
	"os"
	"release-notifier/api"
	"release-notifier/cache"
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
	closer.Bind(main)
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

	closer.Close()
}

func notifyIfNeeded(repo infrastructure.RepositoryConfig) {
	fmt.Println("Checking new version for:", repo.Name)
	var (
		githubToken = os.Getenv("GITHUB_TOKEN")
		releases    = api.GetReleases(repo.Url, githubToken)
	)
	if releases == nil || len(releases) == 0 {
		fmt.Println(fmt.Sprintf("%s: not found releases. Skip", repo.Name))
		return
	}

	var latestRelease = releases[0]

	if isAvailableNewVersion(latestRelease) {
		notify(latestRelease, repo.Name)
	}
}

func isAvailableNewVersion(release api.Release) bool {
	if cache.IsExists(release) == false {
		cache.Save(release)
		return false
	}

	version := cache.GetCurrentVersion(release)
	var isAvailable = version != release.TagName

	if isAvailable {
		fmt.Println("New version available!", version)
	} else {
		fmt.Println("New version not available. Current version", version)
	}

	return isAvailable
}

func notify(release api.Release, name string) {
	var (
		chatId   = infrastructure.GetConfig().TelegramConfig.ChatID
		text     = infrastructure.GetConfig().TelegramConfig.NewVersionMessage
		message  = fmt.Sprintf(text, name, release.TagName, release.Url)
		botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	)

	telegram.SendMessage(chatId, message, botToken)
	cache.Save(release)

	fmt.Println("Save new version to cache", release.TagName)
}
