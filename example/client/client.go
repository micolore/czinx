package main

import (
	"czinx/znet"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {

	fmt.Println("client test start ... ")
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start error exit", err)
		return
	}
	for {
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewPackage(101, []byte("Zinx V0.5 Client Test Message")))
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("write error err", err)
			return
		}

		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println(" read head error ")
			break
		}

		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println(" server unpack data err ", err)
				return
			}
			fmt.Println("===>recv Msg: ID=", msg.Id, ",len= ", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep((1 * time.Second))
	}
}
