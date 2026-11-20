package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// latestURL is GitHub's record of the newest published release: a release that
// is neither a draft nor a prerelease, which is what installScript installs.
const latestURL = "https://api.github.com/repos/bezmoradi/gopamin/releases/latest"

// installScript is the release's one-command installer; it always installs the
// newest published release.
const installScript = "https://github.com/bezmoradi/gopamin/releases/latest/download/install.sh"

// updateWait bounds the whole update check, so it never holds up `gopamin new`.
const updateWait = 3 * time.Second

// updateHint tells how to update to latest: the install script, pointed at the
// folder this gopamin runs from, or `go install` for those who installed with Go.
func updateHint(latest string) string {
	return fmt.Sprintf("\nGopamin %v is available (you have %v). To update:\n  curl -fsSL %v | %vsh\nOr, with Go: go install github.com/bezmoradi/gopamin@%v\nSet GOPAMIN_NO_UPDATE_CHECK=1 to stop this check.\n\n",
		latest, VERSION, installScript, installDirPrefix(), latest)
}

// installDirPrefix returns "GOPAMIN_INSTALL_DIR='<folder>' " for the folder of
// the running gopamin, symlinks resolved and single-quoted for sh, or "" when
// that folder can't be found or holds a character a terminal could misshow.
func installDirPrefix() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return ""
	}
	dir := filepath.Dir(exe)
	if !filepath.IsAbs(dir) || strings.IndexFunc(dir, unsafeRune) >= 0 || !utf8.ValidString(dir) {
		return ""
	}
	return "GOPAMIN_INSTALL_DIR='" + strings.ReplaceAll(dir, "'", `'\''`) + "' "
}

// unsafeRune reports control and invisible formatting characters, such as
// bidirectional overrides, that would hide what a printed command does.
func unsafeRune(r rune) bool {
	return unicode.IsControl(r) || unicode.Is(unicode.Cf, r) || unicode.Is(unicode.Zl, r) || unicode.Is(unicode.Zp, r)
}

// newerVersion returns the latest published release's tag when it is newer
// than VERSION, and "" otherwise. It is advisory: a network failure, a timeout, a
// non-200 answer, a draft or prerelease, or a tag that isn't vMAJOR.MINOR.PATCH
// returns "" within updateWait, so it never blocks a command. Setting
// GOPAMIN_NO_UPDATE_CHECK, CI or GITHUB_ACTIONS turns it off.
func newerVersion() string {
	for _, env := range []string{"GOPAMIN_NO_UPDATE_CHECK", "CI", "GITHUB_ACTIONS"} {
		if os.Getenv(env) != "" {
			return ""
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), updateWait)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, latestURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "gopamin/"+VERSION)
	client := http.Client{Timeout: updateWait}
	res, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		return ""
	}

	var info struct {
		TagName    string `json:"tag_name"`
		Draft      bool   `json:"draft"`
		Prerelease bool   `json:"prerelease"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&info); err != nil || info.Draft || info.Prerelease {
		return ""
	}
	if !strings.HasPrefix(info.TagName, "v") || !isNewer(info.TagName, VERSION) {
		return ""
	}
	return info.TagName
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
		if part == "" || strings.Trim(part, "0123456789") != "" {
			return nil, false
		}
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, false
		}
		numbers[i] = n
	}
	return numbers, true
}
