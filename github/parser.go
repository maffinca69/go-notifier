package github

import "regexp"

func ParseRepositoryUrl(url string) (repo string, name string) {
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
