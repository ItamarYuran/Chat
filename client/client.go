package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	Addr     string
	Lang     string
	Username string
}

func NewClient(addr string) *Client {
	return &Client{Addr: addr}
}

func (c *Client) Start() {
	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		fmt.Println("error connecting so server: ", err)
		return
	}
	defer conn.Close()

	go listenForMessages(conn)

	c.getUsername(conn)
	c.getLanguage(conn)
	c.sendAddr(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("You: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "exit" {
			fmt.Println("tchussi")
			return
		}
		conn.Write([]byte(message + "\n"))
	}
}

func listenForMessages(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected from server")
			return
		}
		fmt.Print(message)
	}

}

func (c *Client) getUsername(conn net.Conn) {
	fmt.Println("Enter Username: ")
	username := bufio.NewReader(os.Stdin)
	name, _ := username.ReadString('\n')
	c.Username = name
	conn.Write(([]byte(name + "\n")))
}

func (c *Client) getLanguage(conn net.Conn) {
	fmt.Println("Enter Language: ")
	language := bufio.NewReader(os.Stdin)
	lang, _ := language.ReadString('\n')
	c.Lang = lang
	conn.Write(([]byte(lang + "\n")))
}

func (c *Client) sendAddr(conn net.Conn) {
	conn.Write(([]byte(c.Addr + "\n")))
}
