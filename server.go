package main

import (
	"fmt"
	"net"
)

// 构建服务端，封装结构体类

type Server struct {
	IP   string // 服务器IP
	Port int    // 服务器端口
}

// 用于创建一个新的Server实例
// 接收ip和port参数， 返回一个Server指针
func NewServer(ip string, port int) *Server {
	return &Server{
		IP:   ip,
		Port: port,
	}
}

// 启动服务器的方法
func (this *Server) Start() {
	// socket listen
	// 创建一个tcp的socket
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.IP, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	// close listener
	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}
		// do handler
		go this.Handler(conn)
	}
}

// 处理连接的业务
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("连接成功")
}
