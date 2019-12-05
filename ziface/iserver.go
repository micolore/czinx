package ziface

//定义服务器相关接口
type Iserver interface {

	//开启服务
	Start()

	//停止服务
	Stop()

	//服务启动
	Serve()

	//v0.6 add msgId
	AddRouter(msgId uint32, router IRouter)

	//获取连接管理器
	GetConnMgr() IConnmanager

	//在连接开始到时候执行
	CallOnConnStart(conn Iconnection)

	//在连接停止到时候执行
	CallOnConnStop(conn Iconnection)

	//设置连接开始的时候
	SetOnConnStart(func(Iconnection))

	//设置连接停止的时候
	SetOnConnStop(func(Iconnection))
}
