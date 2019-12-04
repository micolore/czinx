package znet

import (
	"czinx/ziface"
	"errors"
	"fmt"
	"io"
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

	Router ziface.IRouter

	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {

	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		//管道存储的是协程间通信的相关信息
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader() {

	fmt.Println("reader goroutine is running")

	defer fmt.Println(c.RemoteAddr().String(), "conn reader is exit!")
	defer c.Stop()

	for {

		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTcpConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			continue
		}

		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitBuffChan <- true
			continue
		}

		if msgHead.GetDataLen() > 0 {

			data := make([]byte, msgHead.GetDataLen())
			if _, err := io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitBuffChan <- true
				continue
			}

			msgHead.SetData(data)

		}
		//dataStr := string(msgHead.GetData())
		//fmt.Println("send router data ", dataStr)
		//fmt.Println(data)
		//fmt.Println("conn reev data length ", msg.GetDataLen())
		//fmt.Println("msg str data", str)
		//fmt.Println("msg data = ", data)
		//str := string(data)
		//fmt.Println("msg id = ", msgHead.GetMsgId())
		//fmt.Println("msg len = ", msgHead.GetDataLen())

		req := Request{
			conn: c,
			msg:  msgHead,
		}

		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

func (c *Connection) Start() {

	//开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()
	//开启链接之后开始阻塞，直到链接相关处理执行完成
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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {

	if c.isClosed == true {
		return errors.New("Connection  closed when send msg")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewPackage(msgId, data))
	if err != nil {
		fmt.Println(" pack error msg id = ", msgId)
		return errors.New("pack error msg ")
	}
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("write msg id ", msgId, " error ")
		c.ExitBuffChan <- true
		return errors.New("conn write error")
	}
	return nil
}
