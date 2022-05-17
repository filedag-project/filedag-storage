package mutcask

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"testing"
)

func TestServer(t *testing.T) {
	MutServer("127.0.0.1", "9011", utils.TmpDirPath(t))
}

func TestServer2(t *testing.T) {
	MutServer("127.0.0.1", "9012", utils.TmpDirPath(t))
}

func TestServer3(t *testing.T) {
	MutServer("127.0.0.1", "9013", utils.TmpDirPath(t))
}
func TestHeartBeating(t *testing.T) {
	sendHeartBeat("127.0.0.1:7373")
}
