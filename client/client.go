package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIP   string   // 服务器IP
	ServerPort int      // 服务器端口
	Name       string   // 客户端用户名
	conn       net.Conn // 套接字句柄
	choice     int      // 选择的模式
}

// 创建一个新的客户端实例
func NewClient(serverIP string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		choice:     999, // 默认值
	}

	// 连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn // 绑定连接句柄
	// 返回客户端对象
	return client
}

var serverIp string
var serverPort int

// .\client.exe -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认值:127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口号(默认值:8888)")
}
func (client *Client) memu() bool {
	// 打印菜单
	fmt.Println("请输入你的选择:")
	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新名称")
	fmt.Println("0. 退出聊天")
	var choice int
	fmt.Scanln(&choice)
	if choice >= 0 && choice <= 3 {
		client.choice = choice
		return true
	} else {
		fmt.Println("请输入合法范围内的数字")
		return false
	}
}
func (client *Client) UpdateName() bool {
	// 更新名称
	fmt.Println("请输入新的用户名:")
	fmt.Scanln(&client.Name)
	// 发送给服务器
	sendMsg := "rename:" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	fmt.Println("更新名称成功")
	return true
}
func (client *Client) PublicChat() { // 公聊模式
	// 提示用户输出信息
	var chatMsg string
	fmt.Println("请输入消息内容，以exit退出")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		// 发送消息给服务器
		if len(chatMsg) != 0 { // 消息非空才能发送
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println("请输入消息内容，以exit退出")
		fmt.Scanln(&chatMsg)
	}

}
func (client *Client) Run() {
	for client.choice != 0 {
		for client.memu() == false {
		}

		switch client.choice {
		case 1: // 公聊模式
			client.PublicChat()
			break
		case 2: // 私聊模式
			fmt.Println("私聊模式")
			break
		case 3: // 更新名称
			client.UpdateName()
			break
		}
	}
}

// 接收处理服务器的消息进行终端打印
func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn) // 监听服务器的消息拷贝打印到标准输出，是下面的代码的简化版
	// for {
	// 	buf := make([]byte, 4096)
	// 	client.conn.Read(buf)
	// 	fmt.Println(buf)
	// }
}
func main() {
	// 解析命令行参数
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("=====连接服务器失败=====")
		return
	}
	go client.DealResponse() // 启动协程监听服务器的回执消息
	fmt.Println("=====连接服务器成功=====")
	client.Run()
}
