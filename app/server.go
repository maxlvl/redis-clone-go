package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
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
	reader := bufio.NewReader(conn)
	for {
		command, err := parseCommand(reader)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Error reading  from client: %s", err.Error())
				os.Exit(1)
			}
		}

		response := handleCommand(command)
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing response: %s", err.Error())
			os.Exit(1)
		}
	}
}

func parseCommand(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(line, "*") {
		return nil, fmt.Errorf("Invalid command format")
	}

	numArgs, err := strconv.Atoi(strings.TrimPrefix(strings.TrimSpace(line), "*"))
	if err != nil {
    fmt.Printf("Error parsing number of commands: %s\n", err.Error())
		return nil, err
	}

	args := make([]string, numArgs)
	for i := 0; i < numArgs; i++ {
		line, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if !strings.HasPrefix(line, "$") {
			return nil, fmt.Errorf("invalid argument format")
		}

		length, err := strconv.Atoi(strings.TrimPrefix(strings.TrimSpace(line), "$"))
		if err != nil {
      fmt.Printf("Error parsing length of each command: %s\n", err.Error())
			return nil, err
		}

		arg, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		args[i] = strings.TrimSpace(arg)[:length]
	}

	return args, nil
}

func handleCommand(command []string) string {
	if len(command) == 0 {
		return "-ERR empty command \r\n"
	}

	switch strings.ToUpper(command[0]) {
	case "PING":
		return "+PONG\r\n"
	case "ECHO":
		if len(command) < 2 {
			return "-ERR wrong number of arguments for 'ECHO' command\r\n"
		}

		return fmt.Sprintf("$%d\r\n%s\r\n", len(command[1]), command[1])

	default:
		return fmt.Sprintf("-ERR unknown command '%s'\r\n", command[0])
	}
}
