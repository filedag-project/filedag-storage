package utils

import (
	"io/ioutil"
	"testing"
)

func TmpDirPath(t *testing.T) string {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to make temp dir: %s", err)
	}
	return tmpdir
}
