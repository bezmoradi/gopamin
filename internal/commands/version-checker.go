package commands

import (
	"io"
	"net/http"
	"regexp"
)

const ENV_URL = "https://raw.githubusercontent.com/gopamin/gopamin/master/internal/commands/constants.go"

func versionChecker() (bool, string) {
	res, err := http.Get(ENV_URL)
	if err != nil {
		return true, ""
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return true, ""
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return true, ""
	}

	newVersion := extractVersion(string(bytes))
	if VERSION != newVersion {
		return false, newVersion
	}

	return true, ""
}

func extractVersion(input string) string {
	re := regexp.MustCompile(`VERSION\s*=\s*"v(\d+\.\d+\.\d+)"`)
	matches := re.FindStringSubmatch(input)
	return "v" + matches[1]
}
