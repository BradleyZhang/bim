package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	OnlineUsers map[string]*User
	ouLock      sync.RWMutex

	BroadcastMsg chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:           ip,
		Port:         port,
		OnlineUsers:  make(map[string]*User),
		BroadcastMsg: make(chan string),
	}
}
func (s *Server) Handler(conn net.Conn) {
	user := NewUser(conn)
	s.ouLock.Lock()
	s.OnlineUsers[user.Name] = user
	s.ouLock.Unlock()

	s.BroadCast(user, "已上线")
}
func (s *Server) HandleBroadcast() {
	for {
		msg := <-s.BroadcastMsg

		s.ouLock.Lock()
		for _, cli := range s.OnlineUsers {
			cli.C <- msg
		}
		s.ouLock.Unlock()
	}
}

func (s *Server) BroadCast(user *User, msg string) {
	var sengMsg strings.Builder
	sengMsg.WriteString("[")
	sengMsg.WriteString(user.Addr)
	sengMsg.WriteString("]")
	sengMsg.WriteString(user.Name)
	sengMsg.WriteString(":")
	sengMsg.WriteString(msg)

	s.BroadcastMsg <- sengMsg.String()
}

// 启动
func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

	// 监听 Broadcast Message
	go s.HandleBroadcast()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//do handler
		go s.Handler(conn)
	}
}
