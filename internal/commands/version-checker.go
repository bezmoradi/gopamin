package commands

import (
	"io"
	"net/http"
	"regexp"
	"time"
)

const ENV_URL = "https://raw.githubusercontent.com/bezmoradi/gopamin/main/internal/commands/constants.go"

func versionChecker() (bool, string) {
	client := http.Client{Timeout: 3 * time.Second}
	res, err := client.Get(ENV_URL)
	if err != nil {
		return true, ""
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		return true, ""
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return true, ""
	}

	newVersion := extractVersion(string(bytes))
	if newVersion == "" {
		return true, ""
	}
	if VERSION != newVersion {
		return false, newVersion
	}

	return true, ""
}

func extractVersion(input string) string {
	re := regexp.MustCompile(`VERSION\s*=\s*"v(\d+\.\d+\.\d+)"`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return ""
	}
	return "v" + matches[1]
}
