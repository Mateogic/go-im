package main

// 主入口
func main() {
	server := NewServer("127.0.0.1", 8888)
	server.Start()
}
