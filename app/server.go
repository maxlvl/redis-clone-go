package main

import (
	"fmt"
	"net"
  "bufio"
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
  scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		request := scanner.Text()

		fmt.Println("Request received:", request)

		response := "+PONG\r\n"
		_, err := conn.Write([]byte(response))
		fmt.Println("Sending response:", response)
		if err != nil {
			fmt.Println("Error occurred in sending reply:", err)
			return
		}
	}
}
