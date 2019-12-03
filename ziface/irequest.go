package ziface

type IRequest interface {
	GetConnection() Iconnection
	GetData() []byte
}
