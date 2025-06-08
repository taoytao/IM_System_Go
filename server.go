package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	IP        string
	Port      int
	OnlineMap map[string]*User //在线用户列表
	mapLock   sync.RWMutex     //读写锁
	Message   chan string      //消息广播的channel
}

// 创建一个Server接口
func NewSever(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 监听Message广播消息channel的goroutine
func (server *Server) ListenMessage() {
	for {
		msg := <-server.Message

		server.mapLock.Lock()
		//将msg发送给全部的在线user
		for _, cli := range server.OnlineMap {
			cli.Channel <- msg
		}

		server.mapLock.Unlock()
	}
}

// 广播消息实现
func (server *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + msg

	server.Message <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	// 用户上线, 将用户加入onlineMap中
	server.mapLock.Lock()
	user := NewUser(conn)
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()

	// 广播当前用户上线消息
	server.Broadcast(user, "已上线")

	// 阻塞
	select {}
}

// 启动服务器的接口
func (server *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.IP, server.Port))
	if err != nil {
		fmt.Println("listen server error:", err)
		return
	}

	// close listen socket
	defer listener.Close()

	// 启动监听Message的goroutine
	go server.ListenMessage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		// do handler
		go server.Handler(conn)
	}

}
