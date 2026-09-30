package main

import (
	"bim/command"
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var (
	serverIp   = flag.String("ip", "127.0.0.1", "服务器IP地址（默认127.0.0.1）")
	serverPort = flag.Int("port", 8888, "服务器端口（默认8888）")
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	flag       int //当前模式
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       -1,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn
	return client
}
func (c *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊")
	fmt.Println("2.私聊")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println(">>>>> 请输入合法范围内数字 <<<<<")
		return false
	}
}
func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {

		}
		switch c.flag {
		case 1:
			c.PublishChat()
		case 2:
			c.PrivateChat()
		case 3:
			c.Updatename()
		}
	}
}
func (c *Client) DealResponse() {
	io.Copy(os.Stdout, c.conn)
}
func (c *Client) Updatename() bool {
	fmt.Println(">>>>> 请输入用户名:")
	fmt.Scanln(&c.Name)

	sendMsg := command.Rename + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}
func (c *Client) PublishChat() {
	reader := bufio.NewReader(os.Stdin)
	var chat string
	fmt.Println(">>>>> 请输入聊天内容,exit退出.")
	chat, err := reader.ReadString('\n')
	chat = strings.TrimSpace(chat)

	if err != nil {
		fmt.Println("read err", err)
		return
	}
	for chat != "exit" {
		if len(chat) != 0 {
			sendMsg := chat + "\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}
		chat = ""
		fmt.Println(">>>>> 请输入聊天内容,exit退出.")
		chat, err = reader.ReadString('\n')
		chat = strings.TrimSpace(chat)
		if err != nil {
			fmt.Println("read err", err)
			break
		}
	}
}

func (c *Client) ListUsers() {
	sendMsg := command.GetOnlineUsers + "\n"
	if _, err := c.conn.Write([]byte(sendMsg)); err != nil {
		fmt.Println("conn Write err:", err)
	}
}

func (c *Client) PrivateChat() {
	for {
		c.ListUsers()
		fmt.Println(">>>>> 输入聊天对象用户名,exit 退出")
		reader := bufio.NewReader(os.Stdin)
		remoteUser, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read err:", err)
			break
		}
		remoteUser = strings.TrimSpace(remoteUser)
		if remoteUser == "exit" {
			break
		}
		for {
			fmt.Println(">>>>> 输入聊天信息,exit 退出")
			msg, err := reader.ReadString('\n')
			msg = strings.TrimSpace(msg)
			if err != nil {
				fmt.Println("read err", err)
				break
			}
			if msg == "exit" {
				break
			}
			if len(msg) != 0 {
				sendMsg := command.PrivateChat + remoteUser + command.Delimiter + msg + "\n"
				_, err := c.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

		}
	}
}
func main() {
	flag.Parse()

	client := NewClient(*serverIp, *serverPort)
	if client == nil {
		fmt.Println(">>>>> 连接服务器失败...")
		return
	}
	go client.DealResponse()
	fmt.Println(">>>>> 连接服务器成功...")
	client.Run()
}
