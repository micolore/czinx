package znet

import (
	"czinx/utils"
	"czinx/ziface"
	"fmt"
	"strconv"
)

//定义消息处理
type MsgHandle struct {
	Apis         map[uint32]ziface.IRouter
	WorkPoolSize uint32
	TaskQueue    []chan ziface.IRequest
}

//创建一个消息处理
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:         make(map[uint32]ziface.IRouter),
		WorkPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:    make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

//处理请求
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {

	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println(" Api msgId = ", request.GetMsgID(), " is not found!")
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//添加路由
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {

	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println(" Add Api MsgID = ", msgId)
}

//开启worker
func (mh *MsgHandle) StartOneWorker(workerID int, TaskQueue chan ziface.IRequest) {

	fmt.Println("Worker ID =", workerID, " is Started...")
	for {
		select {
		case request := <-TaskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//初始化工作池,根据系统默认配置
func (mh *MsgHandle) StartWorkerPool() {

	totalSize := int(mh.WorkPoolSize)
	for i := 0; i < totalSize; i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		fmt.Printf("Init worker pool total size is:%d , now init is:%d\n", totalSize, i+1)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}

}

//发送消息到任务队列
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {

	workerID := request.GetConnection().GetConnID() % mh.WorkPoolSize
	fmt.Println(" add connID=", request.GetConnection().GetConnID(), " request msgID =", request.GetMsgID(), "to workerID=", workerID)
	mh.TaskQueue[workerID] <- request

}
