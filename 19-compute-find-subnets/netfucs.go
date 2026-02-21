package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Convert a dots-and-numbers IP address to a single 32-bit numeric
func Ipv4ToValue(addr string) (uint32, error) {
	parts := strings.Split(addr, ".")
	if len(parts) != 4 {
		return 0, fmt.Errorf("invalid IP format")
	}

	var value uint32
	for i := range 4 {
		// Parse directly to uint64 to avoid signed int issues
		num, err := strconv.ParseUint(parts[i], 10, 8)
		if err != nil {
			return 0, fmt.Errorf("invalid octet: %s", parts[i])
		}

		// Shift the byte into the correct position Part 0 shifts 24
		// bits, Part 1 shifts 16, etc.
		value |= uint32(num) << (8 * (3 - i))
	}
	return value, nil
}

// NOTE:
//  << (Left Shift): This is the standard way to move bits.
//  Ex: num << 8 is exactly the same as num * 256.
//  |= (Bitwise OR): Used to combine the shifted octets into the final
//  uint32.

// Convert a single 32-bit numeric value of integer type to a
// dots-and-numbers IP address.
func ValueToIpv4(val uint32) string {
	// Extract each byte by shifting it into the lowest position and
	// masking with 0xFF (binary 11111111)
	octet1 := byte(val >> 24)
	octet2 := byte(val >> 16 & 0xFF)
	octet3 := byte(val >> 8 & 0xFF)
	octet4 := byte(val & 0xFF)
	return fmt.Sprintf("%d.%d.%d.%d", octet1, octet2, octet3, octet4)
}

// NOTE: Why use & 0xFF?
//  If you shift 192.168.0.1 right by 16 bits, you are left with the
//  bits for 192 and 168 combined. The & 0xFF (which is 255 in
//  decimal) acts like a filter that only lets the rightmost 8 bits
//  through, effectively "erasing" the 192 part so you only see the
//  168.

// Given a subnet mask in slash notation, return mask as a single
// uint32 number
func GetSubnetMask(slash string) uint32 {
	_, maskStr, found := strings.Cut(slash, "/")
	if !found {
		return 0
	}
	maskBits, _ := strconv.Atoi(maskStr)
	// handle the edge case for /0
	if maskBits == 0 {
		return 0
	}
	// Use 0xFFFFFFFF (all 1s) and shift the zeros in from the right
	// We use ^uint32(0) to represent 32 bits of 1s safely.
	return ^uint32(0) << (32 - maskBits)
}

// Logic: "IP Address AND Subnet Mask = Network ID"
func IPsOnSameSubnet(ip1, ip2, slash string) (bool, error) {
	value1, err := Ipv4ToValue(ip1)
	if err != nil {
		return false, err
	}
	value2, err := Ipv4ToValue(ip2)
	if err != nil {
		return false, err
	}
	mask := GetSubnetMask(slash)
	return (value1 & mask) == (value2 & mask), nil
}

// GetNetwork returns the network portion of an IPv4 address value.
// Logic: "Network Portion = IP Value AND Netmask"
func GetNetwork(ipValue, netmask uint32) uint32 {
	return ipValue & netmask
}

// Search a dictionary(map) of routers (keyed by router IP) to find
// which router belongs to the same subnet as the given IP.
// Return None if no routers is on the same subnet as the given IP.
func FindRouterForIP(routers map[string]string, target string) (string, error) {
	// Iterate thorugh the map of routers.  routerIP is the key (e.g.,
	// "1.2.3.1") and mask is the value (e.g., "/24")
	for routerIP, mask := range routers {
		same, err := IPsOnSameSubnet(routerIP, target, mask)
		switch {
		case err != nil:
			continue
		case same:
			return routerIP, nil
		}
	}
	return "", nil
}
