package main

import (
	"fmt"
	"net"
)

func main() {
	listener, _ := net.Listen("tcp", ":6379")

	conn, _ := listener.Accept()

	defer conn.Close()

	for {
		response := []byte("+PONG\r\n")
    _, err := conn.Write(response)
		if err != nil {
			fmt.Println("Error occurred in sending reply: %s", err)
      return
		}
	}
}
