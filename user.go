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
	CommandGetOnlineUsers = "list"
	CommandRename         = "rename|"
	CommandPrivateChat    = "to|"
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

	if msg == CommandGetOnlineUsers {
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

		return
	}
	if strings.HasPrefix(msg, CommandRename) {
		newName := strings.TrimPrefix(msg, CommandRename)

		if _, ok := u.server.OnlineUsers[newName]; ok {
			u.SendMsg("当前用户名已被使用")
		} else {
			u.server.ouLock.Lock()
			delete(u.server.OnlineUsers, u.Name)
			u.server.OnlineUsers[newName] = u
			u.server.ouLock.Unlock()
			u.Name = newName
		}

		return
	}
	if strings.HasPrefix(msg, CommandPrivateChat) {

		s := strings.Split(strings.TrimPrefix(msg, CommandPrivateChat), "|")
		remoteName, content := s[0], strings.Join(s[1:], "")
		if content == "" {
			u.SendMsg("未声明消息\n用法: \n\tto|username|message")
			return
		}
		if remoteUser, ok := u.server.OnlineUsers[remoteName]; ok {
			var sengMsg strings.Builder
			sengMsg.WriteString("[")
			sengMsg.WriteString(u.Addr)
			sengMsg.WriteString("]")
			sengMsg.WriteString(u.Name)
			sengMsg.WriteString("(私聊):")
			sengMsg.WriteString(content)
			sengMsg.WriteString("\n")
			remoteUser.SendMsg(sengMsg.String())
			return
		} else {
			u.SendMsg("用户:" + remoteName + " 未上线")
			return
		}
	}

	u.server.BroadCastUserMsg(u, msg)
}
