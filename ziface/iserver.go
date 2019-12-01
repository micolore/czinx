package ziface

//定义服务器接口
type Iserver interface {
	Start()

	Stop()

	Serve()
}
