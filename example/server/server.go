package main

import (
	"czinx/ziface"
	"czinx/znet"
	"fmt"
)

type PingRouter struct {
	znet.BaseRoute
}

func (this *PingRouter) PreHandle(req ziface.IRequest) {

	fmt.Println(" call router PreHandle")

	_, err := req.GetConnection().GetTcpConnection().Write([]byte("before ping ...\n"))
	if err != nil {
		fmt.Println(" call back ping ping ping  error")
	}
}

func (this *PingRouter) Handle(req ziface.IRequest) {

	fmt.Println("call router Handle")
	_, err := req.GetConnection().GetTcpConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println(" call back ping ping ping error")
	}

}

func (this *PingRouter) PostHandle(req ziface.IRequest) {
	fmt.Println("call router PostHandle")
	_, err := req.GetConnection().GetTcpConnection().Write([]byte("after...ping...ping\n"))
	if err != nil {
		fmt.Println(" call back ping ping ping error")
	}
}

func main() {

	s := znet.NewServer("[zinx v-0.3]")

	s.AddRouter(&PingRouter{})

	s.Serve()
}
