package ziface

type IDataPack interface {

	//获取head长度
	GetHeadLen() uint32

	//封包
	Pack(ms IMessage) ([]byte, error)

	//拆包
	Unpack([]byte) (IMessage, error)
}
