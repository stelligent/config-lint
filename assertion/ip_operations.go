package assertion

import (
	"fmt"
	"net"
	"strings"
)

var rfc1918PrivateCIDRs = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}

func getIPObject(addressString string) (net.IP, error) {
	if !strings.Contains(addressString, "/") {
		addressString = fmt.Sprintf("%s/32", addressString)
	}
	ipAddress, _, err := net.ParseCIDR(addressString)
	if err != nil {
		return nil, err
	}
	return ipAddress, nil
}

func isSubnet(ipAddressStr string, supernet string) bool {
	ipAddress, parseError := getIPObject(ipAddressStr)
	if parseError != nil {
		Debugf("%v", parseError)
		return false
	}
	_, superNetwork, err := net.ParseCIDR(supernet)
	if err != nil {
		Debugf("error parsing supernet: %v", err)
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
