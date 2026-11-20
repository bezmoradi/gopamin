package commands

import (
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

// latestURL is the Go module proxy's record of the newest tagged release: what
// `go install github.com/bezmoradi/gopamin@latest` installs.
const latestURL = "https://proxy.golang.org/github.com/bezmoradi/gopamin/@latest"

// installScript is the release's one-command installer; it always installs the
// newest published release.
const installScript = "https://github.com/bezmoradi/gopamin/releases/latest/download/install.sh"

// updateHint tells how to update to latest: the install script, pointed at the
// folder this gopamin runs from, or `go install` for those who installed with Go.
func updateHint(latest string) string {
	return fmt.Sprintf("Gopamin %v is available (you have %v). To update, run:\n  curl -fsSL %v | %vsh\nOr, if you installed it with Go: go install github.com/bezmoradi/gopamin@%v\n\n",
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
