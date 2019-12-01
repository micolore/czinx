package znet

import (
	"czinx/ziface"
	"fmt"
	"net"
)

type Connection struct {
	Conn *net.TCPConn

	//当前连接的id，也可以当作全局的SessionId ID全局唯一
	ConnID uint32

	//当前连接是关闭状态
	isClosed bool

	//该链接的处理方法api
	handleAPI ziface.HandFunc

	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandFunc) *Connection {

	c := &Connection{

		Conn:         conn,
		ConnID:       connID,
		handleAPI:    callback_api,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader() {

	fmt.Println("reader goroutine is running")

	defer fmt.Println(c.RemoteAddr().String(), "conn reader is exit!")
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf error ", err)
			c.ExitBuffChan <- true
			continue
		}
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("connID ", c.ConnID, "handleis is error")
			c.ExitBuffChan <- true
			return
		}
	}
}

func (c *Connection) Start() {

	//开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()

	for {
		select {
		case <-c.ExitBuffChan:
			// 得到退出消息true 不再阻塞
			return
		}
	}
}

func (c *Connection) Stop() {

	if c.isClosed {
		return
	}
	c.isClosed = true

	c.Conn.Close()

	c.ExitBuffChan <- true

	close(c.ExitBuffChan)
}

func (c *Connection) GetTcpConnection() *net.TCPConn {

	return c.Conn
}

func (c *Connection) GetConnID() uint32 {

	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {

	return c.Conn.RemoteAddr()
}
