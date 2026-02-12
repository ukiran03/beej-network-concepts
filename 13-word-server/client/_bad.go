package main

import (
	"encoding/binary"
	"log"
	"net"
)

// This inital method is not Idiomatic

var packetBuf []byte // global

func getNextWordPacketBAD(c net.Conn) []byte {
	res := make([]byte, 0)
	n, err := c.Read(packetBuf)
	if err != nil {
		log.Print(err)
		return nil
	}
	var wlen int
	if n > wordLenSize {
		sizeBytes := packetBuf[:wordLenSize]
		res = append(res, sizeBytes...)                // write the word length size
		wlen = int(binary.BigEndian.Uint16(sizeBytes)) // word length size in int
	}

	if n >= wlen+wordLenSize {
		// write the word itself
		res = append(res, packetBuf[wordLenSize:wordLenSize+wlen]...)
	}
	packetBuf = packetBuf[wordLenSize+wlen:]
	return res
}
