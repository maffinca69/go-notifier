package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Release struct {
	TagName       string `json:"tag_name"`
	Url           string `json:"html_url"`
	RepositoryUrl string
}

const APIReleasesUrl = "https://api.github.com/repos/%s/%s/releases/latest"

func GetLatestRelease(repositoryUrl string) *Release {
	var repo, owner string = ParseRepositoryUrl(repositoryUrl)

	var httpUrl = fmt.Sprintf(APIReleasesUrl, repo, owner)
	return sendRequest(httpUrl)
}

func sendRequest(repositoryUrl string) (release *Release) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", repositoryUrl, nil)

	if err != nil {
		panic("Error creating HTTP request")
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		panic("Error sending HTTP request")
	}

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("Error reading HTTP response body")
	}

	err = json.Unmarshal(body, &release)

	return release
}
