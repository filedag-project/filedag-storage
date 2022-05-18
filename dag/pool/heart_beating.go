package pool

import (
	"net"
	"time"
)

func (r *NodeRecordSys) HandleConnection(conn net.Conn, name string) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			r.Remove(name)
			log.Errorf("%v, Fatal error:%v ", conn.RemoteAddr().String(), err)
			return
		}

		Data := buffer[:n]
		message := make(chan byte)
		go r.HeartBeating(conn, message, 5, name)
		go GravelChannel(Data, message)

		log.Infof("%v,%v,%v", time.Now().Format("2006-01-02 15:04:05.0000000"), conn.RemoteAddr().String(), string(buffer[:n]))
	}
}
func GravelChannel(bytes []byte, mess chan byte) {
	for _, v := range bytes {
		mess <- v
	}
	close(mess)
}
func (r *NodeRecordSys) HeartBeating(conn net.Conn, bytes chan byte, timeout int, name string) {
	select {
	case fk := <-bytes:
		log.Infof("%v heartbeat: the %v times", conn.RemoteAddr().String(), string(fk))
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		break

	case <-time.After(7 * time.Second):
		r.Remove(name)
		log.Errorf("conn dead now")
		conn.Close()
	}
}
