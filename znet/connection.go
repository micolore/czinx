package znet

import (
	"czinx/utils"
	"czinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Connection struct {

	//定义服务器实体
	TcpServer ziface.Iserver

	//定义连接实体
	Conn *net.TCPConn

	//当前连接的id，也可以当作全局的SessionId ID全局唯一
	ConnID uint32

	//当前连接是关闭状态
	isClosed bool

	//该链接的处理方法api
	handleAPI ziface.HandFunc

	//消息处理器
	MsgHandle ziface.IMsgHandle

	//告知该链接已经退出/停止的channel（管道）
	ExitBuffChan chan bool

	//消息channel（管道）
	msgChan chan []byte

	//带buff的channel（管道）
	msgBuffChan chan []byte

	//属性
	property map[string]interface{}

	//锁
	propertyLock sync.RWMutex
}

//创建connection
func NewConnection(server ziface.Iserver, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {

	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		MsgHandle: msgHandler,
		isClosed:  false,
		//管道存储的是协程间通信的相关信息
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte),
		property:     make(map[string]interface{}),
	}
	//v0.8-将新创建的Conn添加到链接管理中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

//开始读
func (c *Connection) StartReader() {

	fmt.Println("reader goroutine is running")

	defer fmt.Println(c.RemoteAddr().String(), "conn reader is exit!")
	defer c.Stop()

	for {

		//创建封/拆包器
		dp := NewDataPack()
		//定义一个slince，长度是dp的head长度
		headData := make([]byte, dp.GetHeadLen())
		//从连接读到slince里面
		if _, err := io.ReadFull(c.GetTcpConnection(), headData); err != nil {
			fmt.Println(" Connection read  head data error", err)
			break
		}
		//拆包headdata
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitBuffChan <- true
			break
		}
		var data []byte
		if msgHead.GetDataLen() > 0 {
			//根据获取的内容长度
			data = make([]byte, msgHead.GetDataLen())
			if _, err := io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Println("Connection read msg data error ", err)
				break
			}
		}
		//把内容设置到消息实体里面
		msgHead.SetData(data)

		//dataStr := string(msgHead.GetData())
		//fmt.Println("send router data ", dataStr)
		//fmt.Println(data)
		//fmt.Println("conn reev data length ", msg.GetDataLen())
		//fmt.Println("msg str data", str)
		//fmt.Println("msg data = ", data)
		//str := string(data)
		//fmt.Println("msg id = ", msgHead.GetMsgId())
		//fmt.Println("msg len = ", msgHead.GetDataLen())

		//创建request到实体
		req := Request{
			conn: c,
			msg:  msgHead,
		}
		//如果系统的默认工作长度满足条件，就使用msghandle发送信息
		if utils.GlobalObject.WorkerPoolSize > 0 {
			go c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			//使用默认的
			go c.MsgHandle.DoMsgHandler(&req)

		}

	}
}

//写消息goroutine
func (c *Connection) StartWriter() {

	fmt.Println(c.RemoteAddr().String(), "[conn write exit!]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer exit!]")

	for {
		//根据具体的channel里面的数据，来决定具体使用哪个回写消息
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error: ,", err, "conn writer exit")
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("send buff data error ,", err, " conn writer exit")
					return
				}
			} else {
				break
				fmt.Println(" msgBuffChan is Closed")
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

//连接开启
func (c *Connection) Start() {

	//开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()

	//开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()

	//调用服务器设置的开启连接方法
	c.TcpServer.CallOnConnStart(c)

	//开启链接之后开始阻塞，直到链接相关处理执行完成
	for {
		select {
		case <-c.ExitBuffChan:
			// 得到退出消息true 不再阻塞
			return
		}
	}
}

//连接关闭
func (c *Connection) Stop() {

	//已经关闭就直接返回
	if c.isClosed {
		return
	}
	//设置连接是否关闭标示为true
	c.isClosed = true

	//调用服务器设置的关闭连接方法
	c.TcpServer.CallOnConnStop(c)

	//连接关闭
	c.Conn.Close()

	//写退出channel消息为true
	c.ExitBuffChan <- true

	//从连接管理器里面移除该连接
	c.TcpServer.GetConnMgr().Remove(c)

	//关闭channel
	close(c.ExitBuffChan)
	close(c.msgChan)
}

//获取tcp connection
func (c *Connection) GetTcpConnection() *net.TCPConn {

	return c.Conn
}

//获取connID
func (c *Connection) GetConnID() uint32 {

	return c.ConnID
}

//获取远程服务地址信息
func (c *Connection) RemoteAddr() net.Addr {

	return c.Conn.RemoteAddr()
}

//发送消息
func (c *Connection) SendMsg(msgId uint32, data []byte) error {

	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewPackage(msgId, data))
	if err != nil {
		fmt.Println(" pack error msg id = ", msgId)
		return errors.New("pack error msg ")
	}

	c.msgChan <- msg

	return nil
}

//发送带缓冲的消息
func (c *Connection) SendBuffMsg(connID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send buff msg")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewPackage(connID, data))

	if err != nil {
		fmt.Println("pack error msg id = ", connID)
		return errors.New("pack error msg ")
	}
	c.msgBuffChan <- msg
	return nil
}

//设置自定义属性
func (c *Connection) SetProperty(key string, value interface{}) {

	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	fmt.Println("connection set property key:", key)
	c.property[key] = value
}

//获取自定义属性
func (c *Connection) GetProperty(key string) (interface{}, error) {

	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除自定义属性
func (c *Connection) RemoveProperty(key string) {

	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
