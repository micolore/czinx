package znet

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

//封包
func NewPackage(id uint32, data []byte) *Message {
	// new 根据位置初始化、根据名称初始化、以及匿名字段的数据类型问题
	// make（map) init的问题
	return &Message{
		Id:      id,
		Data:    data,
		DataLen: uint32(len(data)),
	}
}

//获取数据长度
func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

//获取消息id
func (msg *Message) GetMsgId() uint32 {
	return msg.Id
}

//获取数据
func (msg *Message) GetData() []byte {
	return msg.Data
}

//设置消息数据段长度
func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}

//设计消息ID
func (msg *Message) SetMsgId(msgId uint32) {
	msg.Id = msgId
}

//设计消息内容
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
