package configs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	t.Run("loads .env from the working directory", func(t *testing.T) {
		// The test brings its own .env: the project's is gitignored, so a clean
		// checkout (CI) has none.
		const key = "CONFIGS_TEST_LOADED"
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(key+"=yes\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		t.Chdir(dir)
		t.Cleanup(func() { _ = os.Unsetenv(key) })

		LoadEnv()

		if got := os.Getenv(key); got != "yes" {
			t.Fatalf("expected %s=yes but got %q", key, got)
		}
	})
}
