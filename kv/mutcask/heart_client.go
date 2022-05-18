package mutcask

import (
	"net"
	"os"
	"time"
)

func SendHeartBeat(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		//Log(os.Stderr, "Fatal error:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Errorf("Fatal error:%v", err.Error())
		os.Exit(1)
	}
	log.Infof("%v,connection succcess!", conn.RemoteAddr().String())

	sender(conn)
	log.Infof("send over")
}
func sender(conn *net.TCPConn) {
	for {
		words := "heart beat client."
		msg, err := conn.Write([]byte(words))
		if err != nil {
			log.Errorf(conn.RemoteAddr().String(), "Fatal error: ", err)
			os.Exit(1)
		}
		log.Infof("sever accept:%v", msg)
		time.Sleep(3 * time.Second)
	}
}
