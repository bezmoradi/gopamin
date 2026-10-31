package commands

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

// latestURL is the Go module proxy's record of the newest tagged release: what
// `go install github.com/bezmoradi/gopamin@latest` installs.
const latestURL = "https://proxy.golang.org/github.com/bezmoradi/gopamin/@latest"

// newerVersion returns the latest release when it is newer than VERSION, and ""
// otherwise. It is advisory: a network or parse failure returns "" within 3 s, so
// it never blocks a command. GOPAMIN_NO_UPDATE_CHECK=1 turns it off.
func newerVersion() string {
	if os.Getenv("GOPAMIN_NO_UPDATE_CHECK") != "" {
		return ""
	}

	client := http.Client{Timeout: 3 * time.Second}
	res, err := client.Get(latestURL)
	if err != nil {
		return ""
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		return ""
	}

	var info struct{ Version string }
	if err := json.NewDecoder(io.LimitReader(res.Body, 1<<16)).Decode(&info); err != nil {
		return ""
	}
	if !isNewer(info.Version, VERSION) {
		return ""
	}
	return info.Version
}

// isNewer reports whether version a is newer than b, both vMAJOR.MINOR.PATCH. A
// version it can't parse, such as a pre-release, is never newer.
func isNewer(a, b string) bool {
	versionA, okA := parseVersion(a)
	versionB, okB := parseVersion(b)
	return okA && okB && slices.Compare(versionA, versionB) > 0
}

func parseVersion(version string) ([]int, bool) {
	parts := strings.Split(strings.TrimPrefix(version, "v"), ".")
	if len(parts) != 3 {
		return nil, false
	}
	numbers := make([]int, 3)
	for i, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil || n < 0 {
			return nil, false
		}
		numbers[i] = n
	}
	return numbers, true
}
