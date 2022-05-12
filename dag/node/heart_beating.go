package node

import (
	"fmt"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " Fatal error: ", err)
			return
		}

		Data := buffer[:n]
		message := make(chan byte)
		//心跳计时
		go HeartBeating(conn, message, 4)
		//检测每次是否有数据传入
		go GravelChannel(Data, message)

		Log(time.Now().Format("2006-01-02 15:04:05.0000000"), conn.RemoteAddr().String(), string(buffer[:n]))
	}
}
func GravelChannel(bytes []byte, mess chan byte) {
	for _, v := range bytes {
		mess <- v
	}
	close(mess)
}
func HeartBeating(conn net.Conn, bytes chan byte, timeout int) {
	select {
	case fk := <-bytes:
		Log(conn.RemoteAddr().String(), "心跳:第", string(fk), "times")
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		break

	case <-time.After(5 * time.Second):
		Log("conn dead now")
		conn.Close()
	}
}
func Log(v ...interface{}) {
	fmt.Println(v...)
	return
}
