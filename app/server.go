package main

import (
	"fmt"
	"net"
)

func main() {
	listener, _ := net.Listen("tcp", ":6379")
	defer listener.Close()
	for {
		conn, _ := listener.Accept()
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	response := []byte("+PONG\r\n")
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("Error occurred in sending reply: %s", err)
		return
	}
}
