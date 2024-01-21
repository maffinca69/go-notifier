package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Release struct {
	TagName string `json:"tag_name"`
	Url     string `json:"html_url"`
}

const APIReleasesUrl = "https://api.github.com/repos/%s/%s/releases"

func GetReleases(htmlUrl string, accessToken string) []Release {
	var repo, owner string = parseHtmlUrl(htmlUrl)

	var httpUrl = fmt.Sprintf(APIReleasesUrl, repo, owner)
	return sendRequest(httpUrl, accessToken)
}

func sendRequest(httpUrl string, accessToken string) []Release {
	client := &http.Client{}
	req, err := http.NewRequest("GET", httpUrl, nil)
	if err != nil {
		panic("Error creating HTTP request")
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
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

	releases := make([]Release, 0)

	json.Unmarshal(body, &releases)

	return releases
}

func parseHtmlUrl(url string) (string, string) {
	pattern := regexp.MustCompile("^https?.+(www\\.)?github.com/(?P<repo>[\\w.-]+)/(?P<name>[\\w\\-+]+)")
	match := pattern.FindStringSubmatch(url)
	result := make(map[string]string)

	for i, name := range pattern.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return result["repo"], result["name"]
}