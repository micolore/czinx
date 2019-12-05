package ziface

type IMsgHandle interface {

	//处理请求
	DoMsgHandler(request IRequest)

	//新增router
	AddRouter(msgId uint32, router IRouter)

	//开启一个工作连接池
	StartWorkerPool()

	//发送信息到消息队列
	SendMsgToTaskQueue(rquest IRequest)
}
