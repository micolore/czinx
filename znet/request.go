package znet

import (
	"czinx/ziface"
)

//定义request实体
type Request struct {
	conn ziface.Iconnection
	msg  ziface.IMessage
}

//获取连接信息
func (r *Request) GetConnection() ziface.Iconnection {

	return r.conn
}

//获取内容
func (r *Request) GetData() []byte {

	return r.msg.GetData()
}

//获取消息id
func (r *Request) GetMsgID() uint32 {

	return r.msg.GetMsgId()
}
