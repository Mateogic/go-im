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

// 处理用户的消息
func (this *User) DoMsg(msg string) {
	if msg == "who" { // 查询在线用户
		// 查询在线用户
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			// 发送在线用户列表
			onlineMsg := "[" + user.Addr + "]" + user.Name + ": 在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename:" { // 修改用户名
		newName := msg[7:]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.Conn.Write([]byte("用户名已存在\n"))
		} else {
			// 修改用户名
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.Conn.Write([]byte("用户名成功修改为" + newName + "\n"))
		}
	} else {
		this.server.BroadCast(this, msg)
	}
}

// 查询在线用户
func (this *User) SendMsg(msg string) {
	// 打印在线用户
	this.Conn.Write([]byte(msg))
}
