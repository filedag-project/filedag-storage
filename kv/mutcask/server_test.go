package mutcask

import (
	"testing"
)

func TestServer(t *testing.T) {
	MutServer("127.0.0.1", "9010", "/tmp/dag/data1")
	//go MutServer("127.0.0.1", "9011", utils.TmpDirPath(t))
	//go MutServer("127.0.0.1", "9012", utils.TmpDirPath(t))
}

func TestServer2(t *testing.T) {
	MutServer("127.0.0.1", "9011", "/tmp/dag/data2")
}

func TestServer3(t *testing.T) {
	MutServer("127.0.0.1", "9012", "/tmp/dag/data3")
}

func TestServer4(t *testing.T) {
	MutServer("127.0.0.1", "9013", "/tmp/dag/data4")
}
