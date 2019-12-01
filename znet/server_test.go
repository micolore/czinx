package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

// 如果不想使用test
// 可以单独建两个程序（服务器、客户端），分别进行运行
func ClientTest() {

	fmt.Println("client test ... start")

	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err ", err)
		return
	}
	for {
		_, err := conn.Write([]byte("hello zinx "))
		if err != nil {
			fmt.Println("client write error", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("buf read err", err)
			return
		}
		fmt.Printf("server call back : %s cnt = %d \n", buf, cnt)
		time.Sleep(1 * time.Second)
	}
}

func TestServert(*testing.T) {

	s := NewServer("zinx-v0.1")

	go ClientTest()

	s.Serve()

}
