package main

import "net"

type Server struct {
	IP string
	Port int
}

// 创建一个Server接口
func NewSever(ip string, port int) *Server {
	server := &Server{
        IP: ip,
        Port: port,
    }

	return server
}

func (this *Server) Handler(conn net.Conn) {
	// ...当前连接业务
	fmt.Println("连接建立成功")
	return
}

// 启动服务器的接口
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.IP, this.Port))
	if err != nil {
        fmt.Println("listen server error:", err)
        return
    }

	// close listen socket
	defer listener.Close()

	for {
		// accept
        conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
	}

	// do handler
	go this.Handler(conn)


}