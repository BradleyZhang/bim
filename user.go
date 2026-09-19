package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

const (
	GetOnlineUsers = "list"
)

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user

}

// 监听 chann ，发消息给client
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}

func (u *User) Online() {
	u.server.ouLock.Lock()
	u.server.OnlineUsers[u.Name] = u
	u.server.ouLock.Unlock()

	// 广播上线信息
	u.server.BroadCastUserMsg(u, "已上线")

}
func (u *User) OffLine() {
	u.server.ouLock.Lock()
	delete(u.server.OnlineUsers, u.Name)
	u.server.ouLock.Unlock()

	u.server.BroadCastUserMsg(u, "下线")
}
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg + "\n"))
}
func (u *User) DoMessage(msg string) {
	switch msg {
	case GetOnlineUsers:
		u.server.ouLock.Lock()
		var sengMsg strings.Builder
		for _, user := range u.server.OnlineUsers {
			sengMsg.WriteString("[")
			sengMsg.WriteString(user.Addr)
			sengMsg.WriteString("]")
			sengMsg.WriteString(user.Name)
			sengMsg.WriteString(":在线\n")
		}
		u.SendMsg(sengMsg.String())
		u.server.ouLock.Unlock()
	default:
		u.server.BroadCastUserMsg(u, msg)
	}

}
