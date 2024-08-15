package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	Addr string
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

	fmt.Println("Enter Username: ")
	username := bufio.NewReader(os.Stdin)
	name, _ := username.ReadString('\n')
	conn.Write(([]byte(name)))

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
