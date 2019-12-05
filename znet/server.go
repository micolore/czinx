package znet

import (
	"czinx/utils"
	"czinx/ziface"
	"errors"
	"fmt"
	"net"
	"time"
)

//定义server实体
type Server struct {
	Name      string
	IPversion string
	IP        string
	Port      int
	//怎么理解？
	//Router ziface.IRouter
	msgHandle   ziface.IMsgHandle
	ConnMgr     ziface.IConnmanager
	OnConnStart func(conn ziface.Iconnection)
	OnConnStop  func(conn ziface.Iconnection)
}

//创建服务器实例
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

//服务器启动
func (s *Server) Start() {

	//打印服务器基本配置
	fmt.Printf("[Start] Server Listener at Ip: %s ,Port: %d is starting\n", s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	go func() {

		//初始化worker
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

		fmt.Println("Start Server ", s.Name, " Success! , Now Listenning...")

		var cid uint32
		cid = 0

		// 启动网络连接业务(阻塞)
		for {
			conn, err := listenter.AcceptTCP()
			if err != nil {
				fmt.Println("accept err ", err)
				continue
			}

			//如果连接管理的长度大于系统配置的长度，关闭连接
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

//停止服务器
func (s *Server) Stop() {
	fmt.Println("[Stop] Zinx server ", s.Name)
	s.ConnMgr.CleanConn()
}

//服务器启动服务
func (s *Server) Serve() {
	s.Start()
	for {
		time.Sleep(10 * time.Second)
	}
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

//添加路由
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandle.AddRouter(msgId, router)
	fmt.Println("add router success!")
}

//获取服务器连接管理器
func (s *Server) GetConnMgr() ziface.IConnmanager {
	return s.ConnMgr
}

//设置连接启动的执行方法
func (s *Server) SetOnConnStart(hookFunc func(ziface.Iconnection)) {
	s.OnConnStart = hookFunc
}

//设置连接停止的执行方法
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
