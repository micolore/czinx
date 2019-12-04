package znet

import (
	"czinx/utils"
	"czinx/ziface"
	"fmt"
	"strconv"
)

type MsgHandle struct {
	Apis         map[uint32]ziface.IRouter
	WorkPoolSize uint32
	TaskQueue    []chan ziface.IRequest
}

func NewMsgHandle() *MsgHandle {

	return &MsgHandle{
		Apis:         make(map[uint32]ziface.IRouter),
		WorkPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:    make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {

	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println(" api msgId = ", request.GetMsgID(), " is not found!")
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {

	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println(" add api msgID = ", msgId)
}

func (mh *MsgHandle) StartOneWorker(workerID int, TaskQueue chan ziface.IRequest) {

	fmt.Println("worker id = ", workerID, " is started.")

	for {

		select {
		case request := <-TaskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) StartWorkerPool() {

	for i := 0; i < int(mh.WorkPoolSize); i++ {

		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {

	workerID := request.GetConnection().GetConnID() % mh.WorkPoolSize
	fmt.Println(" add connID=", request.GetConnection().GetConnID(), " request msgID =", request.GetMsgID(), "to workerID=", workerID)
	mh.TaskQueue[workerID] <- request
}
