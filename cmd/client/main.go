package main

import (
	client "Chat/client"
)

func main() {
	cl := client.NewClient("localhost:8080")
	cl.Start()
}
