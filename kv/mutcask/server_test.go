package mutcask

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"os"
	"testing"
)

func TestServer(t *testing.T) {
	os.Setenv(Host, "127.0.0.1")
	os.Setenv(Port, "9011")
	os.Setenv(Path, utils.TmpDirPath(t))
	MutServer()
}

func TestServer2(t *testing.T) {
	os.Setenv(Host, "127.0.0.1")
	os.Setenv(Port, "9012")
	os.Setenv(Path, utils.TmpDirPath(t))
	MutServer()
}

func TestServer3(t *testing.T) {
	os.Setenv(Host, "127.0.0.1")
	os.Setenv(Port, "9013")
	os.Setenv(Path, utils.TmpDirPath(t))
	MutServer()
}
