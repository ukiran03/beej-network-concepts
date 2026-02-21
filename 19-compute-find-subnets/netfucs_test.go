package main

import (
	"fmt"
	"log"
	"testing"
)

// Ipv4ToValue and ValueToIpv4
func TestIpv4RoundTrip(t *testing.T) {
	testcases := []struct {
		addr  string
		value uint32
	}{
		{"255.255.0.0", 4294901760},
		{"1.2.3.4", 16909060},
	}
	for _, tc := range testcases {
		t.Run(tc.addr, func(t *testing.T) {
			valueGot, err := Ipv4ToValue(tc.addr)
			if err != nil {
				t.Fatalf("Failed to convert %s to value: %v", tc.addr, err)
			}
			if valueGot != tc.value {
				t.Errorf("Value mismatch! Got %d, want %d", valueGot, tc.value)
			}

			addrGot := ValueToIpv4(tc.value)
			if addrGot != tc.addr {
				t.Errorf("String mismatch! Got %s, want %s", addrGot, tc.addr)
			}
		})
	}
}

func TestGetSubnetMask(t *testing.T) {
	testcases := []struct {
		slash string
		bits  uint32
	}{
		{"/16", 4294901760},
		{"10.20.30.40/23", 4294966784},
		{"/0", 0},
	}
	for _, tc := range testcases {
		t.Run(tc.slash, func(t *testing.T) {
			result := GetSubnetMask(tc.slash)
			if result != tc.bits {
				t.Errorf("Mismatch! Got %d, want %d", result, tc.bits)
			}
		})
	}
}

func TestIPsOnSameSubnet(t *testing.T) {
	testcases := []struct {
		ip1      string
		ip2      string
		subnet   string
		expected bool
	}{
		{"10.23.121.17", "10.23.121.225", "/23", true},
		{"10.23.230.22", "10.24.121.225", "/16", false},
	}
	for i, tc := range testcases {
		testName := fmt.Sprintf("Test%d", i)
		t.Run(testName, func(t *testing.T) {
			got, err := IPsOnSameSubnet(tc.ip1, tc.ip2, tc.subnet)
			if err != nil {
				log.Fatal(err)
			}
			if got != tc.expected {
				t.Errorf("Got %v, Expected %v", got, tc.expected)
			}
		})
	}
}

func TestGetNetwork(t *testing.T) {
	testcases := []struct {
		ipVal    uint32
		netmask  uint32
		expected uint32
	}{
		// 0x01020304  0xffffff00 0x01020300
		{16909060, 4294967040, 16909056},
	}
	for i, tc := range testcases {
		testName := fmt.Sprintf("Test%d", i)
		t.Run(testName, func(t *testing.T) {
			got := GetNetwork(tc.ipVal, tc.netmask)
			if got != tc.expected {
				t.Errorf("Got %v, Expected %v", got, tc.expected)
			}
		})
	}
}

func TestFindRouterForIP(t *testing.T) {
	testcases := []struct {
		routers  map[string]string
		target   string
		expected string
	}{
		{
			routers: map[string]string{
				"1.2.3.1": "/24",
				"1.2.4.1": "/24",
			},
			target:   "1.2.3.5",
			expected: "1.2.3.1",
		},
		{
			routers: map[string]string{
				"1.2.3.1": "/24",
				"1.2.4.1": "/24",
			},
			target:   "1.2.5.6",
			expected: "",
		},
		{
			routers: map[string]string{
				"10.34.166.1": "/24",
				"10.34.52.1":  "/24",
			},
			target:   "10.34.250.1",
			expected: "",
		},
	}
	for i, tc := range testcases {
		testName := fmt.Sprintf("Test%d", i)
		t.Run(testName, func(t *testing.T) {
			got, _ := FindRouterForIP(tc.routers, tc.target)
			if got != tc.expected {
				t.Errorf("Got %v, Expected %v", got, tc.expected)
			}
		})
	}
}
