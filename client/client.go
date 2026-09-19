package main

import (
	"bim/command"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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
	fmt.Println("0.推出")

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
			fmt.Println("公聊")
		case 2:
			fmt.Println("私聊")
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
