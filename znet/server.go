package znet

import (
	"czinx/ziface"
	"fmt"
	"net"
	"time"
)

type Server struct {
	Name      string
	IPversion string
	IP        string
	Port      int
}

func (s *Server) Start() {
	fmt.Printf("[start] server listener at ip: %s ,port: %d is starting\n", s.IP, s.Port)

	go func() {
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

		// 启动网络连接业务(阻塞)
		for {
			conn, err := listenter.AcceptTCP()
			if err != nil {
				fmt.Println("accept err ", err)
				continue
			}
			go func() {
				for {
					//最大512的回显
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf  err ", err)
						continue
					}
					//回显
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err ", err)
						continue
					}

				}
			}()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] zinx server ", s.Name)
}

func (s *Server) Serve() {

	s.Start()

	for {
		time.Sleep(10 * time.Second)
	}
}

func NewServer(name string) ziface.Iserver {
	s := &Server{
		Name:      name,
		IPversion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
