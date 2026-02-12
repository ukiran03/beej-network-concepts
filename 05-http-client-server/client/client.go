package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	host := flag.String("host", "localhost", "address to send request")
	port := flag.String("port", "8080", "port to send request")
	flag.Parse()

	addr := net.JoinHostPort(*host, *port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	defer conn.Close()

	// write request
	request := fmt.Sprintf(
		"GET / HTTP/1.1\r\n"+
			"Host: %s\r\n"+
			"Connection: close\r\n"+
			"Content-Type: text/plain\r\n"+
			"\r\n"+
			"Hello!\r\n",
		*host,
	)
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Fatalf("Error writing to conn: %v", err)
	}

	// signal that we are done sending (Sends FIN)
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}

	// read response
	resp, err := io.ReadAll(conn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", resp)
}
