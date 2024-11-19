package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reqBuffer := make([]byte, 1024)

	n, err := conn.Read(reqBuffer)
	if err != nil {
		fmt.Println("Failed to read the request:", err)
		return
	}

	req := string(reqBuffer[:n])
	sReq := strings.Split(req, "\r\n")

	// fmt.Printf("\n%+q", sReq)

	requestLine := sReq[0]
	headers := getHeaders(sReq[1:])

	fmt.Printf("\n%+q", headers)

	path := strings.Split(requestLine, " ")[1]

	switch {
	case path == "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	case strings.Split(path, "/")[1] == "echo":
		message := strings.Split(path, "/")[2]
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)))
	case strings.Split(path, "/")[1] == "user-agent":
		userAgent := "test"
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)))
	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func getHeaders(sReq []string) (headers []string) {
	for _, line := range sReq {
		if line == "" {
			return
		}

		headers = append(headers, line)
	}

	return
}
