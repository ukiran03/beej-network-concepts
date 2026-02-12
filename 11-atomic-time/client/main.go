package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	host = "time.nist.gov"
	port = "37"
)

const nistUnixOffset int64 = 2208988800 // Jan 1, 1900 to Jan 1, 1970.

func main() {
	addr := net.JoinHostPort(host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	defer conn.Close()

	unixSeconds := time.Now().Unix() + nistUnixOffset

	var rawSeconds uint32
	err = binary.Read(conn, binary.BigEndian, &rawSeconds)
	if err != nil {
		log.Print("Error reading binary data:", err)
		return
	}

	fmt.Printf("Unix Timestamp: %d\n", unixSeconds)
	fmt.Printf("NIST Timestamp: %d\n", rawSeconds)
}
