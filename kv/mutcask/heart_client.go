package mutcask

import (
	beat "github.com/filedag-project/filedag-storage/dag/node/heart_beat"
	"net"
	"os"
	"time"
)

func sendHeartBeat(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		//Log(os.Stderr, "Fatal error:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		beat.Log("Fatal error:", err.Error())
		os.Exit(1)
	}
	beat.Log(conn.RemoteAddr().String(), "connection succcess!")

	sender(conn)
	beat.Log("send over")
}
func sender(conn *net.TCPConn) {
	for {
		words := "heart beat client."
		msg, err := conn.Write([]byte(words))
		if err != nil {
			beat.Log(conn.RemoteAddr().String(), "Fatal error: ", err)
			os.Exit(1)
		}
		beat.Log("sever accept", msg)
		time.Sleep(3 * time.Second)
	}
}
