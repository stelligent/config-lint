package assertion

import (
	"fmt"
	"log"
	"net"
	"strings"
)

var rfc1918PrivateCIDRs = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}

func getIPObject(addressString string) net.IP {
	if !strings.Contains(addressString, "/") {
		addressString = fmt.Sprintf("%s/32", addressString)
	}
	ipAddress, _, err := net.ParseCIDR(addressString)
	if err != nil {
		log.Fatal("error parsing client ip:", err)
	}
	return ipAddress
}

func isSubnet(ipAddressStr string, supernet string) bool {
	ipAddress := getIPObject(ipAddressStr)
	_, superNetwork, err := net.ParseCIDR(supernet)
	if err != nil {
		log.Fatal("error parsing supernet:", err)
	}
	return superNetwork.Contains(ipAddress)
}

func isPrivateIP(ipAddressStr string) bool {
	for _, cidr := range rfc1918PrivateCIDRs {
		if isSubnet(ipAddressStr, cidr) {
			return true
		}
	}
	return false
}
