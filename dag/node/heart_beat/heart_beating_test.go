package heart_beat

import (
	"net"
	"os"
	"testing"
	"time"
)

func TestHeart_beating(t *testing.T) {
	server := ":7373"
	netListen, err := net.Listen("tcp", server)
	if err != nil {
		Log("connect error: ", err)
		os.Exit(1)
	}
	Log("Waiting for Client ...")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			Log(conn.RemoteAddr().String(), "Fatal error: ", err)
			continue
		}

		//设置短连接(10秒)
		conn.SetReadDeadline(time.Now().Add(time.Duration(10) * time.Second))

		Log(conn.RemoteAddr().String(), "connect success!")
		go HandleConnection(conn)
	}
}
