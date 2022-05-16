package mutcask

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"testing"
)

func TestServer(t *testing.T) {
	mutServer("127.0.0.1", "9001", utils.TmpDirPath(t))
}

func TestServer2(t *testing.T) {
	mutServer("127.0.0.1", "9002", utils.TmpDirPath(t))
}

func TestServer3(t *testing.T) {
	mutServer("127.0.0.1", "9003", utils.TmpDirPath(t))
}
