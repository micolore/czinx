package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	fmt.Println("client test start ... ")
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start error exit", err)
		return
	}
	for {

		_, err := conn.Write([]byte("hello zinx 0.2"))
		if err != nil {
			fmt.Println("write error err", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error", err)
			return
		}
		fmt.Printf("server call back : %s,cnt =%d\n", buf, cnt)
		time.Sleep(3 * time.Second)
	}
}
