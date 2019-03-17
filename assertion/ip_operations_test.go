package assertion

import (
	"testing"
)

var ipTests = []struct {
	value          string
	supernet       string
	expectedResult bool
}{
	{"1.1.1.1", "10.0.0.0/8", false},
	{"1.1.1.1/32", "10.0.0.0/8", false},
	{"10.1.0.0/16", "10.0.0.0/8", true},
	{"10.1.1.1/32", "10.0.0.0/8", true},
	{"10.1.1.1", "10.0.0.0/8", true},
}

func TestIsSubnet(t *testing.T) {
	for _, input := range ipTests {
		t.Run(input.value, func(t *testing.T) {
			result := isSubnet(input.value, input.supernet)
			if result != input.expectedResult {
				t.Errorf("got %v, want %v", result, input.expectedResult)
			}
		})
	}
}

var privateIPTests = []struct {
	value          string
	expectedResult bool
}{
	{"1.1.1.1", false},
	{"1.1.1.1/32", false},
	{"10.1.0.0/16", true},
	{"10.1.1.1/32", true},
	{"10.1.1.1", true},
	{"172.16.1.1", true},
	{"172.15.1.1", false},
	{"192.168.1.1", true},
	{"52.1.1.1", false},
}

func TestIsPrivateIp(t *testing.T) {
	for _, input := range ipTests {
		t.Run(input.value, func(t *testing.T) {
			result := isPrivateIP(input.value)
			if result != input.expectedResult {
				t.Errorf("got %v, want %v", result, input.expectedResult)
			}
		})
	}
}
