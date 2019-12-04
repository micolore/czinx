package ziface

type IConnmanager interface {
	Add(conn Iconnection)
	Remove(conn Iconnection)
	Get(conID uint32) (Iconnection, error)
	Len() int
	CleanConn()
}
