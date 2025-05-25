package main

import "net"

type User struct {
	Name   string      // 用户名
	Addr   string      // 客户端的地址
	C      chan string // 与该用户绑定的消息通道
	Conn   net.Conn    // 与该用户绑定的连接
	server *Server     // 当前用户所在的服务器
}

// 创建一个新的用户实例
func NewUser(conn net.Conn, server *Server) *User {
	// 获取客户端的地址
	userAddr := conn.RemoteAddr().String()
	// 创建一个新的用户实例
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		Conn:   conn,
		server: server,
	}
	// 启动监听消息的协程
	go user.ListenMessage()
	return user
}

// 监听当前用户的消息通道，将消息发送给对端用户
func (this *User) ListenMessage() {
	for {
		// 读取消息
		msg := <-this.C
		// 发送消息
		this.Conn.Write([]byte(msg + "\n"))
	}
}

// 用户上线封装
func (this *User) Online() {
	// 将上线用户添加到在线用户列表
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播当前用户的上线消息
	this.server.BroadCast(this, "上线了")
}

// 用户下线封装
func (this *User) Offline() {
	// 删除当前用户的在线记录
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	// 发送下线消息
	this.server.BroadCast(this, "下线了")
}

// 用户发送广播消息封装
func (this *User) DoMsg(msg string) {
	// 发送消息
	this.server.BroadCast(this, msg)
}
