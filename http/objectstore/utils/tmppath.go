package utils

import (
	"os"
	"testing"
)

func TmpDirPath(t *testing.T) string {
	tmpdir, err := os.MkdirTemp("", "fds")
	if err != nil {
		t.Fatalf("failed to make temp dir: %s", err)
	}
	os.RemoveAll(tmpdir)
	return tmpdir
}
