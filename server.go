package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

// 构建服务端，封装结构体类

type Server struct {
	IP   string // 服务器IP
	Port int    // 服务器端口
	// 在线用户列表
	OnlineMap map[string]*User // key: 用户名, value: 用户对象
	mapLock   sync.RWMutex     // 读写锁

	// 用于消息广播的通道
	Message chan string
}

// 用于创建一个新的Server实例
// 接收ip和port参数， 返回一个Server指针
func NewServer(ip string, port int) *Server {
	return &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User), // 初始化在线用户列表
		Message:   make(chan string),      // 初始化消息通道
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

	// 启动监听消息的协程
	go this.ListenMessager()

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
	// fmt.Println("连接成功")
	// 创建一个新的用户实例
	user := NewUser(conn, this)
	user.Online() // 用户上线
	// 监听用户是否活跃的channel
	isLive := make(chan bool)
	// 接收客户端的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				// 用户下线
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn.Read err:", err)
				return
			}
			// 处理用户的消息
			// msg := string(buf[:n-1])                  // 去除\n
			msg := strings.TrimSpace(string(buf[:n])) // <--- 修改后：使用 TrimSpace 去除首尾空白字符
			user.DoMsg(msg)
			isLive <- true // 任意消息表明用户活跃，重置定时器
		}
	}()
	for {
		// 让处理每个客户端连接的 Handler goroutine 在完成初始化工作后继续存活，以维持连接的有效性
		select {
		case <-isLive: // 利用case穿透执行下面的代码更新计时器
		case <-time.After(time.Second * 100):
			{
				// 10秒后自动断开连接
				user.SendMsg("长时间未操作，您已被踢下线")
				// 销毁资源
				close(user.C)
				// 关闭连接
				conn.Close()
				// 退出协程
				return
			}
		}
	}
}

// 发送广播消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ": " + msg
	this.Message <- sendMsg
}

// 监听广播消息
func (this *Server) ListenMessager() {
	for {
		// 读取消息
		msg := <-this.Message
		// 广播消息给所有在线用户
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}
