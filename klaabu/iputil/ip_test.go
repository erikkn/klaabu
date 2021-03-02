package iputil

import (
	"net"
	"testing"
)

func Test_CloneIP(t *testing.T) {
	ip := net.IPv4(1, 2, 3, 4)
	clone := CloneIP(ip)
	if !clone.Equal(ip) {
		t.Fatalf("Source and clone not equal")
	}
	clone[0] = 5
	if clone.Equal(ip) {
		t.Fatalf("Source and clone share backing array")
	}
}

func Test_NextIP(t *testing.T) {
	type Case struct {
		ip   net.IP
		next net.IP
		fail bool
	}

	cases := []Case{
		{net.IPv4(1, 2, 3, 4), net.IPv4(1, 2, 3, 5), false},
		{net.IPv4(1, 2, 3, 255), net.IPv4(1, 2, 4, 0), false},
		{net.IPv4(1, 2, 255, 255), net.IPv4(1, 3, 0, 0), false},
		{net.IPv4(1, 255, 255, 255), net.IPv4(2, 0, 0, 0), false},
		{net.IPv4(255, 255, 255, 255), net.IPv4(0, 0, 0, 0), true},
	}

	for _, c := range cases {
		ip := c.ip.To4()
		next, err := NextIP(ip)
		if err != nil {
			if !c.fail {
				t.Errorf("Unexpected error: %s: %s", c.ip, err)
			}
		} else {
			if c.fail {
				t.Errorf("Expected error, but did not get any: %s: next = %s", c.ip, next)
			} else {
				if !next.Equal(c.next) {
					t.Errorf("Expectation mismatch: %s: expected %s, got %s", c.ip, c.next, next)
				}
			}
		}
	}
}

func Test_PreviousIP(t *testing.T) {
	type Case struct {
		ip   net.IP
		next net.IP
		fail bool
	}

	cases := []Case{
		{net.IPv4(1, 2, 3, 4), net.IPv4(1, 2, 3, 3), false},
		{net.IPv4(1, 2, 3, 0), net.IPv4(1, 2, 2, 255), false},
		{net.IPv4(1, 2, 0, 0), net.IPv4(1, 1, 255, 255), false},
		{net.IPv4(1, 0, 0, 0), net.IPv4(0, 255, 255, 255), false},
		{net.IPv4(0, 0, 0, 0), net.IPv4(0, 0, 0, 0), true},
	}

	for _, c := range cases {
		ip := c.ip.To4()
		next, err := PreviousIP(ip)
		if err != nil {
			if !c.fail {
				t.Errorf("Unexpected error: %s: %s", c.ip, err)
			}
		} else {
			if c.fail {
				t.Errorf("Expected error, but did not get any: %s: next = %s", c.ip, next)
			} else {
				if !next.Equal(c.next) {
					t.Errorf("Expectation mismatch: %s: expected %s, got %s", c.ip, c.next, next)
				}
			}
		}
	}
}
