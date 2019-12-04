package main

import (
	"czinx/ziface"
	"czinx/znet"
	"fmt"
)

type PingRouter struct {
	znet.BaseRoute
}

type HelloZinxRouter struct {
	znet.BaseRoute
}

func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("call helloZinxRouter handle")
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("hello zinx router v0.6"))
	if err != nil {
		fmt.Println(err)
	}
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

func DoConnectionBegin(conn ziface.Iconnection) {

	fmt.Println("DoConnectionBegin is called")

	err := conn.SendMsg(2, []byte("DoConnection Begin..."))

	if err != nil {

		fmt.Println(err)
	}
}

func DoConnectionLost(conn ziface.Iconnection) {

	fmt.Println("DoConnection Lost...")
}

func main() {

	s := znet.NewServer("[zinx v-0.9]")
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	s.Serve()
}
