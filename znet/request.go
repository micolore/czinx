package znet

import (
	"czinx/ziface"
)

type Request struct {
	conn ziface.Iconnection
	data []byte
}

func (r *Request) GetConnection() ziface.Iconnection {

	return r.conn
}

func (r *Request) GetData() []byte {

	return r.data
}
