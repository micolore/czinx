package ziface

type IRequest interface {

	//获取连接
	GetConnection() Iconnection

	//获取数据
	GetData() []byte

	//获取消息id
	GetMsgID() uint32
}
