package main

import (
	client "github.com/ItamarYuran/Chat/client"
)

func main() {
	cl := client.NewClient("localhost:8080")
	cl.Start()
}
