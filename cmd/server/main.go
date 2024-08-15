package main

import (
	server "Chat/server"
)

func main() {
	serv := server.NewServer(":8080")
	serv.Start()
}
