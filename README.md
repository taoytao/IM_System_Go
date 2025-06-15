# IM_System_Go
即时通讯系统_Golang

使用Go语言简易实现了一个即时通讯系统服务端
监听本地地址8888端口

## 编译及运行测试
### 1. 编译
go build -o server main.go server.go
### 2. 服务端运行
./server
### 3. 创建客户端
nc 127.0.0.1 8888
### 4. 用户上线广播功能
### 5. 重命名功能
rename|NEWNAME
### 6. 查询当前在线用户功能
who
### 7. 私聊用户功能, 消息格式: to|张三|消息内容
to|username|msg
