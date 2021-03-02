package klaabu

import (
	"net"
	"reflect"
	"testing"
)

func TestCidr_MinMaxIP(t *testing.T) {
	data := []string{
		"10.0.0.0/8", "10.0.0.0", "10.255.255.255",
		"10.0.0.0/30", "10.0.0.0", "10.0.0.3",
		"10.0.0.0/32", "10.0.0.0", "10.0.0.0",
		"10.1.1.1/32", "10.1.1.1", "10.1.1.1",
		"172.24.32.0/20", "172.24.32.0", "172.24.47.255",
		"0.0.0.0/0", "0.0.0.0", "255.255.255.255",
	}

	for i := 0; i < len(data); i += 3 {
		cidr := Cidr(data[i])
		expectedMin := net.ParseIP(data[i+1])
		expectedMax := net.ParseIP(data[i+2])

		min, max, err := cidr.MinMaxIP()
		if err != nil {
			t.Errorf("error: %s", err)
		}

		if !min.Equal(expectedMin) {
			t.Errorf("Expected: %v, got: %v", expectedMin, min)
		}

		if !max.Equal(expectedMax) {
			t.Errorf("Expected: %v, got: %v", expectedMax, max)
		}
	}
}

func Test_ParseCidr(t *testing.T) {
	_, ipNet, err := net.ParseCIDR("10.0.0.0/8")
	if err != nil {
		t.Errorf("error: %s", err)
	}

	if !reflect.DeepEqual([]byte(ipNet.IP), []byte{10, 0, 0, 0}) {
		t.Errorf("Expected: 10.0.0.0 got: %v", ipNet.IP)
	}

	if !reflect.DeepEqual([]byte(ipNet.Mask), []byte{255, 0, 0, 0}) {
		t.Errorf("Expected: 10.255.255.255, got: %v", ipNet.Mask)
	}
}
