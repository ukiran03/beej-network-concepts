package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const SERVE_FILES = "./testdata"

func main() {
	port := flag.String("port", "28333", "port to listen request")
	flag.Parse()

	addr := ":" + *port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	log.Printf("Server listening at port %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()

	var resp bytes.Buffer
	defer resp.WriteTo(c)

	method, file := parseReq(c)
	if method == "" || file == "" {
		sendError(&resp, "400 Bad Request")
		return
	}

	safePath := filepath.Join(SERVE_FILES, filepath.Clean("/"+file))

	// DEBUG:
	absPath, _ := filepath.Abs(safePath)
	fmt.Printf("Attempting to serve: %s\n", absPath)

	if method == "GET" {
		data, err := os.ReadFile(safePath)
		if err != nil {
			if os.IsNotExist(err) {
				sendError(&resp, "404 File not found")
				return
			} else {
				sendError(&resp, "500 Internal Server Error")
			}
			return
		}

		var ftype string = "text/plain"
		switch filepath.Ext(safePath) {
		case ".html":
			ftype = "text/html"
		case ".jpg", ".jpeg":
			ftype = "image/jpeg"
		}

		buildResp(&resp, "200 OK", ftype, strconv.Itoa(len(data)))
		resp.Write(data)
	}
}

func sendError(b *bytes.Buffer, status string) {
	buildResp(b, status, "text/plain", strconv.Itoa(len(status)))
	b.WriteString(status)
}

func buildResp(b *bytes.Buffer, status, ctype, clen string) {
	// HTTP standards require \r\n (CRLF)
	heads := fmt.Sprintf(
		"HTTP/1.1 %s\r\n"+
			"Content-Type: %s\r\n"+
			"Content-Length: %s\r\n"+
			"Connection: close\r\n\r\n",
		status, ctype, clen)
	b.Write([]byte(heads))
}

// parseReq returns Method (GET) and FilePath (/file1.txt) from the
// request, nothing else.
func parseReq(in io.Reader) (method, file string) {
	scanner := bufio.NewScanner(in)
	if !scanner.Scan() {
		return "", "" //  empty input
	}

	// example: GET /file2.html HTTP/1.1
	header := scanner.Text()
	parts := strings.Fields(header)
	if len(parts) < 2 {
		return "", ""
	}

	method = parts[0]
	file = stripPrefixSlash(parts[1])
	return method, file
}

func stripPrefixSlash(s string) string {
	return strings.TrimPrefix(s, "/")
}
