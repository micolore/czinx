package main

import (
	"czinx/ziface"
	"czinx/znet"
	"fmt"
)

//定义PingRouter
type PingRouter struct {
	znet.BaseRoute
}

//定义HelloZinxRouter
type HelloZinxRouter struct {
	znet.BaseRoute
}

//HelloZinxRouter的实现
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("call helloZinxRouter handle")
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("hello zinx router v0.6"))
	if err != nil {
		fmt.Println(err)
	}
}

//PingRouter的实现
func (this *PingRouter) Handle(req ziface.IRequest) {

	fmt.Println("call router Handle")
	msgDataStr := string(req.GetData())

	fmt.Println("recv from client : msgId=", req.GetMsgID(), ", data=", msgDataStr)
	err := req.GetConnection().SendMsg(1, []byte(" ping...ping...ping... "))

	if err != nil {
		fmt.Println(" call back ping ping ping error")
	}

}

//请求连接开始的时候执行
func DoConnectionBegin(conn ziface.Iconnection) {

	fmt.Println("DoConnectionBegin is called")
	err := conn.SendMsg(2, []byte("DoConnection Begin..."))

	conn.SetProperty("name", "xiaowang")
	conn.SetProperty("home", "helan")

	if err != nil {
		fmt.Println(err)
	}
}

//请求连接关闭的时候执行
func DoConnectionLost(conn ziface.Iconnection) {

	if name, err := conn.GetProperty("name"); err == nil {

		fmt.Println("conn property name= ", name)
	}

	if home, err := conn.GetProperty("home"); err == nil {

		fmt.Println("conn property home= ", home)
	}

	fmt.Println("DoConnection Lost...")
}

func main() {

	//创建server
	s := znet.NewServer("[zinx v-0.9]")

	//设置服务器启动执行的方法
	s.SetOnConnStart(DoConnectionBegin)

	//设置服务器停止执行的方法
	s.SetOnConnStop(DoConnectionLost)

	//添加路由到apis。根据msgID进行区分
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	//开启服务
	s.Serve()
}
