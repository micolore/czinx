package ziface

type IMessage interface {

	//获取数据长度
	GetDataLen() uint32
	//获取信息ID
	GetMsgId() uint32
	//获取数据
	GetData() []byte

	//设置信息id
	SetMsgId(uint32)
	//设置数据
	SetData([]byte)
	//设置数据长度
	SetDataLen(uint32)
}
