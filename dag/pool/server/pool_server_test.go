package server

import (
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	go StartTestDagPoolServer(t)
	time.Sleep(time.Minute)
}
