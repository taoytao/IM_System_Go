package main

import (
	"net"
	"strings"
)

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn
	server  *Server
}

// 创建用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		conn:    conn,
		server:  server,
	}

	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 用户上线业务
func (user *User) Online() {
	// 用户上线
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	// 广播当前用户上线消息
	user.server.Broadcast(user, "已上线")
}

// 用户下线业务
func (user *User) Offline() {
	// 用户下线, 将用户从onlineMap中删除
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	// 广播当前用户上线消息
	user.server.Broadcast(user, "用户下线")
}

// 发送信息
func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

// 用户处理消息业务
func (user *User) MsgProcess(msg string) {
	if msg == "who" {
		// 查询当前在线用户
		user.server.mapLock.Lock()
		for _, usr := range user.server.OnlineMap {
			onlineMsg := "[" + usr.Addr + "]" + usr.Name + "在线\n"
			user.SendMsg(onlineMsg)
		}
		user.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式: rename|newName
		newName := msg[7:]

		// 判断newName是否存在
		_, ok := user.server.OnlineMap[newName]
		if ok {
			user.SendMsg("当前用户名已存在")
		} else {
			user.server.mapLock.Lock()
			oldName := user.Name
			delete(user.server.OnlineMap, oldName)
			user.Name = newName
			user.server.OnlineMap[newName] = user
			user.server.mapLock.Unlock()
			user.SendMsg("成功更新用户名:" + user.Name + "\n")
		}
	} else if len(msg) > 3 && msg[:3] == "to|" {
		// 消息格式: to|张三|消息内容

		// 1. 获取用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			user.SendMsg("消息格式不正确, 请使用 \"to|张三|你好\"格式。\n")
			return
		}

		// 2. 获取到用户名对应的User对象
		remoteUser, ok := user.server.OnlineMap[remoteName]
		if !ok {
			user.SendMsg("该用户名不存在\n")
			return
		}

		// 3. 发送消息
		content := strings.Split(msg, "|")[2]
		if content == "" {
			user.SendMsg("消息内容为空\n")
			return
		}
		remoteUser.SendMsg(user.Name + ":" + content + "\n")

	} else {
		user.server.Broadcast(user, msg)
	}

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
