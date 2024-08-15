package server

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var (
	clients   = make(map[net.Conn]string)
	broadcast = make(chan string)
	mutex     = &sync.Mutex{}
)

type Server struct {
	Addr string
}

func NewServer(addr string) *Server {
	return &Server{Addr: addr}
}

func (s *Server) Start() {
	server, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return
	}
	defer server.Close()

	fmt.Printf("Server is listening on port %s", s.Addr)

	go handleBroadcast()

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error recieving connection: ", err)
			continue
		}
		go handleClient(conn)

	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Enter you username: "))
	username, _ := bufio.NewReader(conn).ReadString('\n')
	username = username[:len(username)-1]

	mutex.Lock()
	clients[conn] = username
	mutex.Unlock()

	broadcast <- fmt.Sprintf("%s has joined the chat", username)

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			broadcast <- fmt.Sprintf("%s has left the chat", username)
			return
		}
		broadcast <- fmt.Sprintf("%s: %s", username, message)
	}

}

func handleBroadcast() {
	for {
		message := <-broadcast
		mutex.Lock()
		for conn := range clients {
			_, err := conn.Write([]byte(message + "\n"))
			if err != nil {
				fmt.Println("Error sending message: ", err)
				conn.Close()
				delete(clients, conn)
			}

		}
		mutex.Unlock()
	}
}
