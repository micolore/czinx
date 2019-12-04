package znet

import (
	"czinx/utils"
	"czinx/ziface"
	"errors"
	"fmt"
	"net"
	"time"
)

type Server struct {
	Name      string
	IPversion string
	IP        string
	Port      int
	//怎么理解？
	//Router ziface.IRouter
	msgHandle ziface.IMsgHandle
	ConnMgr   ziface.IConnmanager

	OnConnStart func(conn ziface.Iconnection)
	OnConnStop  func(conn ziface.Iconnection)
}

func (s *Server) Start() {
	fmt.Printf("[start] server listener at ip: %s ,port: %d is starting\n", s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
	go func() {

		s.msgHandle.StartWorkerPool()

		//获取一个tcp addr
		addr, err := net.ResolveTCPAddr(s.IPversion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolveTcp addr err", err)
			return
		}
		//监听服务器地址
		listenter, err := net.ListenTCP(s.IPversion, addr)
		if err != nil {
			fmt.Println("listen", s.IPversion, "err", err)
			return
		}
		fmt.Println("start zinx server ", s.Name, "success ,now listenning...")

		var cid uint32
		cid = 0

		// 启动网络连接业务(阻塞)
		for {

			conn, err := listenter.AcceptTCP()
			if err != nil {
				fmt.Println("accept err ", err)
				continue
			}

			//v0.8
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			//处理该新链接请求的业务方法，此时应该有handle和conn绑定的
			dealConn := NewConnection(s, conn, cid, s.msgHandle)
			cid++

			//启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server ", s.Name)

	s.ConnMgr.CleanConn()
}

func (s *Server) Serve() {

	s.Start()

	for {
		time.Sleep(10 * time.Second)
	}
}

func NewServer(name string) ziface.Iserver {
	utils.GlobalObject.Reload()

	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPversion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		msgHandle: NewMsgHandle(),
		ConnMgr:   NewConnmanager(),
	}
	return s
}

// 定义当前客户端链接的handle api
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {

	fmt.Println("[Conn Handle] CallBackToClient ... ")
	if _, err := conn.Write(data[:cnt]); err != nil {

		fmt.Println("Write back buf err ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandle.AddRouter(msgId, router)
	fmt.Println("add router success!")
}

func (s *Server) GetConnMgr() ziface.IConnmanager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(hookFunc func(ziface.Iconnection)) {

	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(ziface.Iconnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.Iconnection) {
	if s.OnConnStart != nil {

		fmt.Println("====>callOnConnStart")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.Iconnection) {
	if s.OnConnStop != nil {

		fmt.Println("====>callOnConnStart")
		s.OnConnStop(conn)
	}
}
