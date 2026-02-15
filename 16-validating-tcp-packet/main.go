package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type inputFileSet struct {
	addrFile string
	dataFile string
}

func NewInputFileSet(i int) (inputFileSet, error) {
	addrFile := fmt.Sprintf("testdata/tcp_addrs_%d.txt", i)
	dataFile := fmt.Sprintf("testdata/tcp_data_%d.dat", i)
	if !fileExists(addrFile) || !fileExists(dataFile) {
		return inputFileSet{},
			fmt.Errorf("Error: File set missing for %d", i)
	}
	return inputFileSet{addrFile, dataFile}, nil
}

type tcpAddrs struct {
	srcIP []byte
	dstIP []byte
}

type inputDataSet struct {
	tcpAddrs
	tcpData []byte
	err     error
}

func originalChecksum(in inputDataSet) uint16 {
	return binary.BigEndian.Uint16(in.tcpData[16:18])
}

func calculatedChecksum(in inputDataSet) uint16 {
	zeroChecksumedTCPData := make([]byte, len(in.tcpData))
	copy(zeroChecksumedTCPData, in.tcpData)
	zeroChecksumedTCPData[16], zeroChecksumedTCPData[17] = 0, 0 // zero checksum

	pseudoHeader := ipPseudoHeader(in.srcIP, in.dstIP, uint16(len(in.tcpData)))

	// concatenate pseudoHeader and zeroChecksumedTCPData
	buf := make([]byte, len(pseudoHeader)+len(zeroChecksumedTCPData))
	copy(buf[0:12], pseudoHeader)         // copy pseudoHeader into the start
	copy(buf[12:], zeroChecksumedTCPData) // copy zero checksumed tcp data into rest

	// handle Odd-Byte padding
	if len(buf)%2 != 0 {
		buf = append(buf, 0)
	}
	// compute our Checksum for the data
	return computeChecksum(buf)
}

func openFileSet(set inputFileSet) (inputDataSet, error) {
	var dataSet inputDataSet
	// Addr files
	ipData, err := os.ReadFile(set.addrFile)
	addrs := strings.Fields(string(ipData))
	srcIP, _ := ipAddrToBytes(addrs[0])
	dstIP, _ := ipAddrToBytes(addrs[1])
	dataSet.srcIP, dataSet.dstIP = srcIP, dstIP
	if err != nil {
		return inputDataSet{}, err
	}
	// TCP Data files
	tcpData, err := os.ReadFile(set.dataFile)
	dataSet.tcpData = tcpData
	if err != nil {
		return inputDataSet{}, err
	}
	return dataSet, nil
}

func ipAddrToBytes(addr string) ([]byte, error) {
	parts := strings.Split(addr, ".")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid IP format")
	}
	res := make([]byte, 4) // pre-allocate exactly 4 bytes
	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return nil, fmt.Errorf("invalid octet: %s", part)
		}
		res[i] = byte(num)
	}
	return res, nil
}

func ipPseudoHeader(srcIP, dstIP []byte, tcpLen uint16) []byte {
	// Assume srcIP and dstIP are 4-byte slices
	// Assume tcpLength is an int or uint16
	header := make([]byte, 12)

	copy(header[0:4], srcIP) // Src IP
	copy(header[4:8], dstIP) // Dst IP
	header[8] = 0            // Reserved Zero byte
	header[9] = 6            // Protocol (6 for TCP)
	binary.BigEndian.PutUint16(header[10:12], uint16(tcpLen))

	return header
}

func main() {
	totalFiles := 10      // 0-9
	var wg sync.WaitGroup // 1. Initialize the WaitGroup

	for i := 0; i < totalFiles; i++ {
		wg.Add(1) // 2. Increment the counter

		go func(i int) {
			defer wg.Done()
			fileSet, err := NewInputFileSet(i)
			if err != nil {
				log.Printf("File %d error: %v", i, err)
				return
			}
			dataSet, err := openFileSet(fileSet)
			if err != nil {
				log.Printf("Open error: %v", err)
				return
			}
			origChecksum := originalChecksum(dataSet)
			calcChecksum := calculatedChecksum(dataSet)
			if calcChecksum == origChecksum {
				fmt.Printf("FileSet %d: PASS\n", i)
			} else {
				fmt.Printf(
					"FileSet %d: FAIL (expected %04x, got %04x)\n",
					i, origChecksum, calcChecksum,
				)
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("Processing complete.")
}
