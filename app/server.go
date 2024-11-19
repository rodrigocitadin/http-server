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

	requestLine, headers, body := parseReq(sReq)

	fmt.Printf("\nReq Line: %s", requestLine)
	fmt.Printf("\nHeaders: %+q", headers)
	fmt.Printf("\nBody: %s", body)
	fmt.Print("\r\n\r\n\r\n")

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

func parseReq(sReq []string) (requestLine string, headers []string, body string) {
	requestLine = sReq[0]

	for i, line := range sReq[1:] {
		if line == "" {
			body = strings.Join(sReq[i+1:], " ")
			return
		}

		headers = append(headers, line)
	}
	return
}
