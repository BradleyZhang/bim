package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"
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

const (
	IdleExpireSecond = 120 // 120s 后超时下线
)

func (s *Server) Handler(conn net.Conn) {
	user := NewUser(conn, s)
	user.Online()

	isLive := make(chan bool)

	// 接收消息转交user处理
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.OffLine()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			msg := string(buf[:n-1]) // 去除'\n'
			user.DoMessage(msg)

			isLive <- true
		}
	}()

	for {
		select {
		case <-time.After(time.Second * IdleExpireSecond): // 超时，
			user.SendMsg("超时下线\n")
			close(user.C)
			user.conn.Close()
			runtime.Goexit()
		case <-isLive:
			// 更新计时器
		}
	}
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

func (s *Server) BroadCastUserMsg(user *User, msg string) {
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
