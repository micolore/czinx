package ziface

import (
	"net"
)

//定义连接接口
type Iconnection interface {
	Start()

	Stop()

	//从当前连接获取原始的socket TcpConn
	GetTcpConnection() *net.TCPConn

	//获取连接id
	GetConnID() uint32

	//获取远程客户端地址信息
	RemoteAddr() net.Addr

	//发送消息-no buff
	SendMsg(msgId uint32, data []byte) error

	//发送消息-have buff
	SendBuffMsg(msgID uint32, data []byte) error

	//设置属性
	SetProperty(key string, value interface{})

	//获取自定义属性
	GetProperty(key string) (interface{}, error)

	//删除属性信息
	RemoveProperty(key string)
}

//定义一个统一处理链接业务的接口
//第一个原生的socket链接 第二个参数是客户端请求的数据 第三个是客户端请求的数据长度
type HandFunc func(*net.TCPConn, []byte, int) error
