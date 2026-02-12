package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"strings"
)

// How many bytes is the word length?
const WORD_LEN_SIZE = 2

func buildWordPacket(count int) ([]byte, []string) {
	var packet bytes.Buffer
	var wlist []string

	buf := make([]byte, WORD_LEN_SIZE) // tmp scratch buffer

	for range count {
		word := WORDS[rand.IntN(len(WORDS))]
		wbytes := []byte(word)

		// write uint16 to the scratch buffer
		binary.BigEndian.PutUint16(buf, uint16(len(wbytes)))

		packet.Write(buf)    // write the length bytes to the buffer
		packet.Write(wbytes) // write the actual word bytes

		wlist = append(wlist, word)
	}
	return packet.Bytes(), wlist
}

func sendWords(c net.Conn) ([]string, error) {
	wcount := rand.IntN(9) + 1 // [1, 10)
	packet, wlist := buildWordPacket(wcount)
	_, err := c.Write(packet)
	if err != nil {
		return nil, err
	}
	return wlist, nil
}

func main() {
	port := flag.String("port", "28333", "port to listen request")
	flag.Parse()

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	// this will close the client after the first read
	// defer listener.Close()

	log.Printf("Waiting for connections on port %s", *port)

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
	log.Printf("Got connection from %v", c.RemoteAddr())

	wlist, err := sendWords(c)
	if err != nil {
		log.Printf("Failed to send words to %v: %v", c.RemoteAddr(), err)
		return
	}
	fmt.Printf("Sent words: %s\n", strings.Join(wlist, ","))
}
