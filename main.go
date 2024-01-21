package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"release-notifier/api"
	"release-notifier/cache"
	"release-notifier/telegram"
	"sync"
	"time"
)

type Repository struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type TelegramConfig struct {
	ChatID            int64  `json:"chat_id"`
	NewVersionMessage string `json:"new_version_message"`
}

type Config struct {
	Repository     []Repository   `json:"repositories"`
	TelegramConfig TelegramConfig `json:"telegram"`
}

const ConfigName = "config.json"

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	start := time.Now()
	defer fmt.Println("Execution Time: ", time.Since(start))

	var payload = loadConfig()

	ctx := context.Background()
	var wg sync.WaitGroup

	for _, repo := range payload.Repository {
		wg.Add(1)
		fmt.Println("Checking new version for:", repo.Name)
		repo := repo
		go func() {
			defer wg.Done()
			var githubToken = os.Getenv("GITHUB_TOKEN")
			var releases = api.GetReleases(repo.Url, githubToken)
			if releases == nil || len(releases) == 0 {
				fmt.Println(fmt.Sprintf("%s: not found releases. Skip", repo.Name))
				return
			}

			var latestRelease = releases[0]

			if isAvailableNewVersion(ctx, latestRelease) {
				notify(ctx, payload, latestRelease, repo.Name)
			}
		}()
	}

	wg.Wait()
}

func loadConfig() Config {
	content, err := os.ReadFile(ConfigName)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	payload := Config{}

	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return payload
}

func isAvailableNewVersion(context context.Context, release api.Release) bool {
	var key = getCacheKey(release)
	client := cache.GetClient()
	var exists, _ = client.Exists(context, key).Result()
	if exists == 0 {
		client.Set(context, key, release.TagName, 0)
		fmt.Println("Doesn't exists version. Saved version", release.TagName)
		return false
	}

	version, _ := client.Get(context, key).Result()
	var isAvailable = version != release.TagName
	if isAvailable {
		fmt.Println("New version available!", version)
	} else {
		fmt.Println("New version not available. Current version", version)
	}

	return isAvailable
}

func getCacheKey(release api.Release) string {
	return fmt.Sprintf("%s_%s", release.Url, release.TagName)
}

func notify(context context.Context, config Config, release api.Release, name string) {
	var chatId = config.TelegramConfig.ChatID
	var text = config.TelegramConfig.NewVersionMessage

	var message = fmt.Sprintf(text, name, release.TagName, release.Url)
	var botToken = os.Getenv("TELEGRAM_BOT_TOKEN")

	telegram.SendMessage(chatId, message, botToken)

	var key = getCacheKey(release)

	client := cache.GetClient()
	err := client.Set(context, key, release.TagName, 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Save new version to cache", release.TagName)
}
