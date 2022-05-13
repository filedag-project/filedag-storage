package mutcask

import "testing"

func TestServer(t *testing.T) {
	go mutServer("127.0.0.1", "9001")
	go mutServer("127.0.0.1", "9002")
	go mutServer("127.0.0.1", "9003")
}
