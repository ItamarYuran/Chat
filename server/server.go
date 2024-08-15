package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync"

	"cloud.google.com/go/translate"
	"github.com/ItamarYuran/Chat/client"
	"golang.org/x/text/language"
)

var (
	clients   = make(map[net.Conn]client.Client)
	broadcast = make(chan string)
	mutex     = &sync.Mutex{}
)

type Server struct {
	Addr            string
	TranslateClient *translate.Client
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

	s.getTranslateClient()
	defer s.TranslateClient.Close()

	fmt.Printf("Server is listening on port %s\n", s.Addr)

	go s.handleBroadcast()

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

	username, _ := bufio.NewReader(conn).ReadString('\n')
	username = username[:len(username)-1]
	language, _ := bufio.NewReader(conn).ReadString('\n')
	language = language[:len(language)-1]
	addres, _ := bufio.NewReader(conn).ReadString('\n')
	addres = addres[:len(addres)-1]

	mutex.Lock()
	clients[conn] = client.Client{
		Addr:     addres,
		Lang:     language,
		Username: username,
	}
	mutex.Unlock()

	welcomeNewClient(clients[conn])

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

func (s *Server) handleBroadcast() {
	for {
		message := <-broadcast
		mutex.Lock()
		for conn := range clients {
			trans, err := s.WriteTranslated(conn, message)
			conn.Write([]byte(trans + "\n"))
			if err != nil {
				fmt.Println("Error sending message: ", err)
				conn.Close()
				delete(clients, conn)
			}

		}
		mutex.Unlock()
	}
}

func welcomeNewClient(c client.Client) {
	fmt.Printf("New user joined, %s, he speaks %s, his address is %s\n", c.Username, c.Lang, c.Addr)
	broadcast <- fmt.Sprintf("%s has joined the chat", c.Username)
}

func (s *Server) getTranslateClient() {
	client, _ := NewTranslationClient("../../server/keys.json")
	s.TranslateClient = client
}

func (s *Server) WriteTranslated(conn net.Conn, msg string) (string, error) {
	lang := clients[conn].Lang
	fmt.Println(lang)
	tag, err := language.Parse(lang)
	if err != nil {
		return "", fmt.Errorf("here")
	}
	fmt.Println(tag, lang)

	translatedTexts, err := TranslateText(context.Background(), s.TranslateClient, []string{msg}, tag)
	if err != nil {
		return "", fmt.Errorf("or here")
	}
	return translatedTexts[0], nil
}
