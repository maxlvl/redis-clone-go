package main

import (
	"fmt"
	 "net"
   "bufio"
)

func main() {
  listener, _ := net.Listen("tcp", ":6379")

  conn, _ := listener.Accept()

  for {
    message, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
      fmt.Println("Error occurred: %s", err)
    }
    fmt.Print("Message received:", string(message))
  }
}
