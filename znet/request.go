package znet

import (
	"czinx/ziface"
)

type Request struct {
	conn ziface.Iconnection
	msg  ziface.IMessage
}

func (r *Request) GetConnection() ziface.Iconnection {

	return r.conn
}

func (r *Request) GetData() []byte {

	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {

	return r.msg.GetMsgId()
}
