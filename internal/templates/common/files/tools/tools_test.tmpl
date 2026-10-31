package tools

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRootDir(t *testing.T) {
	t.Run("finds the nearest directory holding go.mod", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		nested := filepath.Join(root, "a", "b")
		if err := os.MkdirAll(nested, 0o755); err != nil {
			t.Fatal(err)
		}
		t.Chdir(nested)

		got, err := FindRootDir()
		if err != nil {
			t.Fatalf("expected %s but got error %v", root, err)
		}
		// The temp dir may sit behind a symlink (macOS /var -> /private/var).
		if !sameDir(t, got, root) {
			t.Fatalf("expected %s but got %s", root, got)
		}
	})

	t.Run("fails outside a module", func(t *testing.T) {
		t.Chdir(t.TempDir())

		if got, err := FindRootDir(); err == nil {
			t.Fatalf("expected an error but got %s", got)
		}
	})
}

func sameDir(t *testing.T, a, b string) bool {
	t.Helper()
	resolvedA, errA := filepath.EvalSymlinks(a)
	resolvedB, errB := filepath.EvalSymlinks(b)
	if errA != nil || errB != nil {
		t.Fatalf("resolve %s, %s: %v, %v", a, b, errA, errB)
	}
	return resolvedA == resolvedB
}
