package znet

import (
	"czinx/ziface"
	"errors"
	"fmt"
	"sync"
)

//定义连接管理器
type ConnManager struct {
	connections map[uint32]ziface.Iconnection
	connLock    sync.RWMutex
}

//创建连接管理器
func NewConnmanager() *ConnManager {

	return &ConnManager{
		connections: make(map[uint32]ziface.Iconnection),
	}
}

//新增连接到连接管理器里面
func (connMgr *ConnManager) Add(conn ziface.Iconnection) {

	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connection add  to connections successfully conn num = ", connMgr.Len())
}

//从连接管理器里面移除连接
func (connMgr *ConnManager) Remove(conn ziface.Iconnection) {

	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("  connect remove connID=", conn.GetConnID(), "successfully conn num ", connMgr.Len())
}

//根据连接id获取连接
func (connMgr *ConnManager) Get(connID uint32) (ziface.Iconnection, error) {

	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")

	}
}

//获取连接管理器的长度（map）
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

//清理连接管理器
func (connMgr *ConnManager) CleanConn() {

	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID, conn := range connMgr.connections {
		//先停止后删除
		conn.Stop()
		delete(connMgr.connections, connID)
	}
	fmt.Println(" Clear All Connections sucessfully : conn num=", connMgr.Len())
}
