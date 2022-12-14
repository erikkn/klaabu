package klaabu

import (
	"fmt"
	"net"

	"github.com/transferwise/klaabu/klaabu/iputil"
)

// MinMaxIP lalala.
func (c *Cidr) MinMaxIP() (net.IP, net.IP, error) {
	_, ipNet, err := net.ParseCIDR(string(*c))
	if err != nil {
		return nil, nil, fmt.Errorf("error while parsing your CIDR %v with error: %s", string(*c), err)
	}

	min := make([]byte, len(ipNet.IP))
	max := make([]byte, len(ipNet.IP))
	for i := range ipNet.IP {
		min[i] = ipNet.Mask[i] & ipNet.IP[i]
		max[i] = ipNet.Mask[i]&ipNet.IP[i] | ^ipNet.Mask[i]
	}

	return min, max, nil
}

// IsChildOf llalala
func (c *Cidr) IsChildOf(parent Cidr) (bool, error) {
	minChild, maxChild, err := c.MinMaxIP()
	if err != nil {
		return false, fmt.Errorf("error: %s", err)
	}

	minParent, maxParent, err := parent.MinMaxIP()
	if err != nil {
		return false, fmt.Errorf("error: %s", err)
	}

	minComparison, err := iputil.CompareIPs(minParent, minChild)
	if err != nil {
		return false, err
	}
	maxComparison, err := iputil.CompareIPs(maxChild, maxParent)
	if err != nil {
		return false, err
	}

	return minComparison <= 0 && maxComparison <= 0, nil
}

// OverlapsCidr checks if `c` and `cidr` are overlapping or not and returns a boolean.
func (c *Cidr) OverlapsCidr(cidr Cidr) (bool, error) {

	minA, maxA, err := c.MinMaxIP()
	if err != nil {
		return false, fmt.Errorf("error while calculating the min & max IPs of %s with error message: %s", *c, err)
	}

	minB, maxB, err := cidr.MinMaxIP()
	if err != nil {
		return false, fmt.Errorf("error while calculating the min & max IPs of %s with error message: %s", cidr, err)

	}

	maxACmpMinB, err := iputil.CompareIPs(maxA, minB)
	if err != nil {
		return false, fmt.Errorf("unable to compare %s and %s with error message: %s", maxA, minB, err)
	}

	minACmpMaxB, err := iputil.CompareIPs(minA, maxB)
	if err != nil {
		return false, fmt.Errorf("unable to compare %s and %s with error message: %s", minA, maxB, err)
	}

	return maxACmpMinB >= 0 && minACmpMaxB <= 0, nil
}

// ValidateChildrenOverlap checks if there are any overlaps between immediate children of a prefix.
func (p *Prefix) ValidateChildrenOverlap() error {
	children := make([]*Prefix, 0, len(p.Children))

	for _, v := range p.Children {
		children = append(children, v)
	}

	for x := 0; x < len(children); x++ {
		for y := x + 1; y < len(children); y++ {
			overlap, err := children[x].Cidr.OverlapsCidr(children[y].Cidr)
			if err != nil {
				return fmt.Errorf("unable to call `OverlapsCidr` with error message: %s", err)
			}

			if overlap {
				return fmt.Errorf("%s is overlapping with %s", children[x].Cidr, children[y].Cidr)
			}
		}
	}

	return nil
}

// Validate validates a CIDR.
func (c *Cidr) Validate() error {
	_, _, err := net.ParseCIDR(string(*c))
	if err != nil {
		return err
	}

	return nil
}

// Validate checks if you are stupid or not.
func (p *Prefix) Validate() error {
	// Checks the vadility of the actual CIDR of the Prefix
	err := p.Cidr.Validate()
	if err != nil {
		return fmt.Errorf("invalid CIDR '%s': %s", p.Cidr, err)
	}

	// Checks if the children has any overlapping CIDR.
	err = p.ValidateChildrenOverlap()
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}

	for _, v := range p.Children {
		isChild, err := v.Cidr.IsChildOf(p.Cidr)
		if err != nil {
			return err
		}

		if !isChild {
			return fmt.Errorf("%v is not a valid child CIDR of %v", v.Cidr, p.Cidr)
		}

		err = v.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
