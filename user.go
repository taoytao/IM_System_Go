package main

import "net"

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn
}

// 创建用户
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		conn:    conn,
	}

	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 监听当前user channel，处理消息
func (user *User) ListenMessage() {
	for {
		msg, ok := <-user.Channel
		if !ok {
			return
		}
		user.conn.Write([]byte(msg + "\n"))
	}
}
