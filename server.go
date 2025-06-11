package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	server.Message <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	user := NewUser(conn, server)

	// 用户上线, 将用户加入onlineMap中
	user.Online()

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			// 存在非法操作
			if err != nil && err != io.EOF {
				fmt.Println("conn.Read error:", err)
				return
			}

			// 提取用户消息(去除'\n')
			msg := string(buf[:n-1])

			// 用户在线则发送消息
			user.MsgProcess(msg)

			// 用户活跃, 重置定时器
			isLive <- true
		}
	}()

	// 设置一个定时器, 超时用户下线
	for {
		select {
		case <-isLive:
			//当前用户活跃, 重置定时器: 不做处理, 激活time.After(time.Second * 10)
		case <-time.After(time.Second * 100):
			// 超时删除用户
			user.SendMsg("你已超时, 被踢出")

			// 销毁用户资源
			close(user.Channel)

			// 关闭连接
			conn.Close()

			return
		}
	}
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
