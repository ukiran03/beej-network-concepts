package main

import (
	"fmt"
	"sort"
	"strconv"
)

// HexToDecimal converts a hex string (e.g., "0x01020300") to a
// uint32.
func HexToDecimal(hexStr string) (uint32, error) {
	// ParseUint returns uint64, so we cast to uint32
	val, err := strconv.ParseUint(hexStr, 0, 32)
	return uint32(val), err
}

// DecimalStringToUint32 converts a decimal string (e.g., "16909056")
// to a uint32.
func DecimalStringToUint32(deciStr string) (uint32, error) {
	val, err := strconv.ParseUint(deciStr, 10, 32)
	return uint32(val), err
}

// Uint32ToHexStr converts the actual integer to a Hexadecimal string
// for display.
func Uint32ToHexStr(val uint32) string {
	return fmt.Sprintf("0x%08x", val)
}

func uniqIPs(srcdstPairs [][]string) []string {
	m := make(map[string]struct{})
	for _, pair := range srcdstPairs {
		m[pair[0]] = struct{}{}
		m[pair[1]] = struct{}{}
	}
	ips := make([]string, 0, len(m))
	for ip := range m {
		ips = append(ips, ip)
	}
	sort.Strings(ips)
	return ips
}
