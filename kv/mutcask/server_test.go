package mutcask

import "testing"

func TestServer(t *testing.T) {
	mutServer("127.0.0.1", "9001", "/tmp/dag/data1")
}

func TestServer2(t *testing.T) {
	mutServer("127.0.0.1", "9002", "/tmp/dag/data2")
}

func TestServer3(t *testing.T) {
	mutServer("127.0.0.1", "9003", "/tmp/dag/data3")
}
