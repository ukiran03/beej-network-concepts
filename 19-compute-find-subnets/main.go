package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

type Connection struct {
	Netmask   string `json:"netmask"`
	Interface string `json:"interface"`
	AD        int    `json:"ad"`
}

type Router struct {
	Connections map[string]Connection `json:"connections"`
	Netmask     string                `json:"netmask"`
	IfCount     int                   `json:"if_count"`
	IfPrefix    string                `json:"if_prefix"`
}

type NetworkConfig struct {
	Routers map[string]Router `json:"routers"`
	SrcDest [][]string        `json:"src-dest"`
}

func main() {
	js, err := os.ReadFile("testdata/data.json")
	if err != nil {
		log.Fatal(err)
	}

	var data NetworkConfig
	err = json.Unmarshal(js, &data)
	if err != nil {
		log.Fatal(err)
	}

	printRouters(data.Routers)
	fmt.Println()
	printSameSubnets(data.SrcDest)
	fmt.Println()
	printIPRouters(data)
}

func printSameSubnets(srcdstPairs [][]string) {
	const mask = "/24"
	fmt.Println("IP Pairs:")
	for _, pair := range srcdstPairs {
		if len(pair) < 2 {
			fmt.Println("Error: invalid pair provided")
			continue
		}
		src, dst := pair[0], pair[1]
		b, err := IPsOnSameSubnet(src, dst, mask)
		if err != nil {
			log.Printf("Error processing %s-%s: %v", src, dst, err)
			continue
		}
		status := "different subnets"
		if b {
			status = "same subnet"
		}

		fmt.Printf(" %5s\t%s: \t%s\n", src, dst, status)
	}
}

func printRouters(routerTable map[string]Router) {
	var netmaskIP, networkIP string
	fmt.Println("Routers:")
	for rip, router := range routerTable {
		for _, conn := range router.Connections {
			// get the netmask
			slash := conn.Netmask
			netmaskValue := GetSubnetMask(slash)
			netmaskIP = ValueToIpv4(netmaskValue)

			// get the network number
			routerValue, _ := Ipv4ToValue(rip)
			networkValue := GetNetwork(routerValue, netmaskValue)
			networkIP = ValueToIpv4(networkValue)
		}
		fmt.Printf(" %5s: \tnetmask %s: \tnetwork %s\n", rip, netmaskIP, networkIP)
	}
}

// Routers and corresponding IPs
func printIPRouters(data NetworkConfig) {
	routers := makeRoutersMap(data.Routers)
	allIPs := uniqIPs(data.SrcDest)
	fmt.Println("Routers and corresponding IPs:")

	// Map Routers to their associated Host IPs
	routerHostMap := make(map[string][]string)

	for _, ip := range allIPs {
		router, err := FindRouterForIP(routers, ip)
		if err != nil || router == "" {
			router = "None"
		}
		routerHostMap[router] = append(routerHostMap[router], ip)
	}

	// Sort the router keys for consistent output
	var sortedRouters []string
	for r := range routerHostMap {
		sortedRouters = append(sortedRouters, r)
	}
	sort.Strings(sortedRouters)

	for _, rIP := range sortedRouters {
		fmt.Printf(" %5s: \t%v\n", rIP, routerHostMap[rIP])
	}
}

func makeRoutersMap(routers map[string]Router) map[string]string {
	routersMap := make(map[string]string, 3*len(routers))
	for _, router := range routers {
		// We pass the same map into each call, to fill it up directly
		collectRouterConnections(router, routersMap)
	}
	return routersMap
}

func collectRouterConnections(router Router, dst map[string]string) {
	for cip, conn := range router.Connections {
		dst[cip] = conn.Netmask
	}
}
