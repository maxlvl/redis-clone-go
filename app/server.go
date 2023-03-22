package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
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
	kvStore := make(map[string]map[string]interface{})
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

		response := handleCommand(command, kvStore)
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

func handleCommand(command []string, kvStore map[string]map[string]interface{}) string {
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
	case "SET":
		// if a PX is sent along with the SET command, set it in a nested hash map of kvStore.
		// {
		//  "value": "stuff",
		//  "time_set": "current_timestamp",
		//  "px": "expiry time stored in miliseconds"
		// }
		// when doing a GET, check for the presence of the PX value in nested hash_map
		// then, compare current_time - px <= time_set, if so, key is good. If not, key has expired and return NULL value
		if len(command) == 3 {
			key_name := command[1]
			value := command[2]
			kvStore[key_name] = map[string]interface{}{
				"value":    value,
				"time_set": time.Now(),
			}
			response := "+OK\r\n"
			return response
		} else if len(command) == 5 {
			key_name := command[1]
			value := command[2]
			px := command[4]
			kvStore[key_name] = map[string]interface{}{
				"value":    value,
				"time_set": time.Now(),
				"px":       px,
			}
			response := "+OK\r\n"
			return response
		} else {
			return "-ERR wrong number of arguments for SET command\r\n"
		}

	case "GET":
		if len(command) == 2 {
			inner_map, ok := kvStore[command[1]]
			if !ok {
				fmt.Println("Something went wrong trying to fetch the innermap")
				return "-ERR something went wrong BLOOP BLURP\r\n"
			}
			value, ok := inner_map["value"]
			if !ok {
				fmt.Println("Something went wrong trying to fetch the value from the inner_map")
				return "-ERR something went wrong BLOOP BLURP\r\n"
			}

			time_set, ok := inner_map["time_set"].(time.Time)
			if !ok {
				fmt.Println("Something went wrong trying to fetch the time_set from the inner_map")
				return "-ERR something went wrong BLOOP BLURP\r\n"
			}

			px, ok := inner_map["px"].(int64)

			if ok {
				current_time := time.Now()
				time_elapsed := current_time.Add(-time.duration(px) * time.Milisecond)
				if time_elapsed.Before(time_set) {
					response := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
					return response
				} else {
					response := "$-1\r\n"
					return response
				}
			}

		} else {
			return "-ERR wrong number of arguments for GET command\r\n"
		}

	default:
		return fmt.Sprintf("-ERR unknown command '%s'\r\n", command[0])
	}
}

