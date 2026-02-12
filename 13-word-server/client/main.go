package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

// How many bytes is the word length?
const wordLenSize = 2

func main() {
	host := flag.String("host", "localhost", "address to send request")
	port := flag.String("port", "28333", "port to send request")
	flag.Parse()

	addr := net.JoinHostPort(*host, *port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	defer conn.Close()
	fmt.Println("Getting words:")

	for {
		wordPacket, err := getNextWordPacket(conn)
		if err != nil || len(wordPacket) == 0 {
			break
		}
		fmt.Println(extractWord(wordPacket))
	}
}

func extractWord(packet []byte) string {
	return string(packet[wordLenSize:])
}

func getNextWordPacket(c net.Conn) ([]byte, error) {
	lenBuf := make([]byte, wordLenSize)
	if _, err := io.ReadFull(c, lenBuf); err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("failed to read length: %w", err)
	}
	wlen := binary.BigEndian.Uint16(lenBuf)

	wordBuf := make([]byte, wlen) // to read exactly 'wlen' bytes
	if _, err := io.ReadFull(c, wordBuf); err != nil {
		return nil, fmt.Errorf("failed to read payload: %w", err)
	}

	return append(lenBuf, wordBuf...), nil
}

// [12-02-2026] NOTE: A TCP connection is like a pipe, not a filing
// cabinet. When you call Read(), you are literally pulling bytes out
// of a buffer managed by the Operating System. Once those bytes are
// read into your application, they are "consumed" and removed from
// the OS receive buffer.
