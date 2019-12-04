package main

import (
	"czinx/ziface"
	"czinx/znet"
	"fmt"
)

type PingRouter struct {
	znet.BaseRoute
}

func (this *PingRouter) Handle(req ziface.IRequest) {

	fmt.Println("call router Handle")
	msgDataStr := string(req.GetData())
	fmt.Println("recv from client : msgId=", req.GetMsgID(), ", data=", msgDataStr)

	err := req.GetConnection().SendMsg(1, []byte(" ping...ping...ping... "))

	if err != nil {
		fmt.Println(" call back ping ping ping error")
	}

}

func main() {

	s := znet.NewServer("[zinx v-0.5]")

	s.AddRouter(&PingRouter{})

	s.Serve()
}
