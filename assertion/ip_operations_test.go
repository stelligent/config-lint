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
	{"172.16.0.0/12", true},
	{"172.0.0.0/8", false},
	{"172.16.1.1", true},
	{"172.15.1.1", false},
	{"192.168.1.1", true},
	{"52.1.1.1", false},
	{"sg-1234567", false},
}

func TestIsPrivateIp(t *testing.T) {
	for _, input := range privateIPTests {
		t.Run(input.value, func(t *testing.T) {
			result := isPrivateIP(input.value)
			if result != input.expectedResult {
				t.Errorf("got %v, want %v", result, input.expectedResult)
			}
		})
	}
}

var maxHostCountTests = []struct {
	value          string
	max            string
	expectedResult bool
}{
	{"10.0.0.0/8", "1000", false},
	{"10.0.0.0/23", "500", false},
	{"10.1.0.0/16", "65600", true},
	{"10.1.1.1/32", "2", true},
	{"10.1.1.1/32", "1", true},
	{"sg-1234567", "0", false},
}

func TestMaxHostCount(t *testing.T) {
	for _, input := range maxHostCountTests {
		t.Run(input.value, func(t *testing.T) {
			result := maxHostCount(input.value, input.max)
			if result != input.expectedResult {
				t.Errorf("got %v, want %v", result, input.expectedResult)
			}
		})
	}
}
