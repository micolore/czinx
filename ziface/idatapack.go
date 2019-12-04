package ziface

type IDataPack interface {
	GetHeadLen() uint32
	Pack(ms IMessage) ([]byte, error)
	Unpack([]byte) (IMessage, error)
}
