package main

import (
	"encoding/binary"
	"errors"
	"os"
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true // File exists
	}
	if errors.Is(err, os.ErrNotExist) {
		return false // File does not exist
	}
	return false
}

// Thanks to Gemini.
func computeChecksum(data []byte) uint16 {
	var sum uint32

	// Iterate in 16-bit (2 byte) steps
	for i := 0; i < len(data)-1; i += 2 {
		// combine two 8-bit bytes into one 16-bit word (Big Endian)
		sum += uint32(binary.BigEndian.Uint16(data[i : i+2]))
	}

	// Fold the 32-bit sum into 16 bits.  While there is a carry (bits
	// above 0xFFFF), add it back to the bottom.
	for (sum >> 16) > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	// One's complement (bit inversion)
	return ^uint16(sum)
}
