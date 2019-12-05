package ziface

type IConnmanager interface {
	//添加一个连接
	Add(conn Iconnection)

	//移除一个连接
	Remove(conn Iconnection)

	//获取一个连接，根据connID
	Get(connID uint32) (Iconnection, error)

	//长度
	Len() int

	//清理连接
	CleanConn()
}
