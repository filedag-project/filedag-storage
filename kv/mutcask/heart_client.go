package mutcask

import (
	beat "github.com/filedag-project/filedag-storage/dag/node/heart_beat"
	"net"
	"os"
	"strconv"
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
	for i := 0; i < 10; i++ {
		words := strconv.Itoa(i) + " Hello I'm MyHeartbeat Client."
		msg, err := conn.Write([]byte(words))
		if err != nil {
			beat.Log(conn.RemoteAddr().String(), "Fatal error: ", err)
			os.Exit(1)
		}
		beat.Log("sever accept", msg)
		time.Sleep(2 * time.Second)
	}
	for i := 0; i < 2; i++ {
		time.Sleep(12 * time.Second)
	}
	for i := 0; i < 10; i++ {
		words := strconv.Itoa(i) + " Hi I'm MyHeartbeat Client."
		msg, err := conn.Write([]byte(words))
		if err != nil {
			beat.Log(conn.RemoteAddr().String(), "Fatal error: ", err)
			os.Exit(1)
		}
		beat.Log("sever accept", msg)
		time.Sleep(2 * time.Second)
	}

}
