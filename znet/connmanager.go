package znet

import (
	"czinx/ziface"
	"errors"
	"fmt"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.Iconnection
	connLock    sync.RWMutex
}

func NewConnmanager() *ConnManager {

	return &ConnManager{

		connections: make(map[uint32]ziface.Iconnection),
	}
}

func (connMgr *ConnManager) Add(conn ziface.Iconnection) {

	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("connection add  to connections successfully conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn ziface.Iconnection) {

	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections, conn.GetConnID())

	fmt.Println("  connect remove connID=", conn.GetConnID(), "successfully conn num ", connMgr.Len())
}

func (connMgr *ConnManager) Get(connID uint32) (ziface.Iconnection, error) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")

	}
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) CleanConn() {

	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID, conn := range connMgr.connections {

		conn.Stop()

		delete(connMgr.connections, connID)
	}
	fmt.Println(" Clear All Connections sucessfully : conn num=", connMgr.Len())
}
