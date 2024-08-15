package main

import (
	server "github.com/ItamarYuran/Chat/server"
)

func main() {
	serv := server.NewServer(":8080")
	serv.Start()
}
