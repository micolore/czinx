package ziface

//定义服务器接口
type Iserver interface {
	Start()

	Stop()

	Serve()

	//主要给客户端链接处理使用
	//v0.6 add msgId
	AddRouter(msgId uint32, router IRouter)

	GetConnMgr() IConnmanager

	CallOnConnStart(conn Iconnection)

	CallOnConnStop(conn Iconnection)

	SetOnConnStart(func(Iconnection))

	SetOnConnStop(func(Iconnection))
}
