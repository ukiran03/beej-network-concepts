package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	port := flag.String("port", "28333", "port to listen request")
	flag.Parse()

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

// simple response for request
func handleConn(c net.Conn) {
	defer c.Close()

	reqMethod, reqBody, err := parseReq(c)
	if err != nil {
		log.Printf("Error parsing Request from %v: %v",
			c.RemoteAddr().String(), err)
		return
	}

	log.Printf("Request from: %v\nMethod: %v\nBody: %v",
		c.RemoteAddr().String(), reqMethod, reqBody)

	payload := "Hello from server!\n"
	resp := fmt.Sprintf(
		"HTTP/1.1 200 OK\n"+
			"Content-Type: text/plain\n"+
			"Content-Length: %v\n"+
			"Connection: close\n"+
			"\n"+
			payload,
		len(payload))

	c.Write([]byte(resp))
}

// parseReq returns the method and payload of the request nothing else
func parseReq(in io.Reader) (string, string, error) {
	reader := bufio.NewReader(in)
	// Parse Request Line (Method)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	method := strings.Fields(line)[0]
	// Skip all Headers
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", "", err
		}
		if line == "\r\n" || line == "\n" {
			break
		}
	}
	body, err := io.ReadAll(reader)
	if err != nil && err != io.EOF {
		return "", "", err
	}
	return method, string(body), nil
}
