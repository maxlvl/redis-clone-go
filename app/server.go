package main

import (
	"fmt"
	"io"
	"net"
	"os"
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
	for {
		buf := make([]byte, 1024)
		if _, err := conn.Read(buf); err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Error reading  from client: %s", err.Error())
				os.Exit(1)
			}
		}
		response := []byte("+PONG\r\n")
		_, err := conn.Write(response)
		if err != nil {
			fmt.Println("Error writing response: %s", err.Error())
			os.Exit(1)
		}

	}

	// Send "+PONG\r\n" response
}
