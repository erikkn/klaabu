package klaabu

import (
	"fmt"
	"github.com/erikkn/klaabu/klaabu/iputil"
	"sort"
)

// Cidr represents a single CIDR notation.
type Cidr string

// Prefix is data for the schema.
type Prefix struct {
	Aliases  []string
	Cidr     Cidr
	Labels   map[string]string
	Parent   *Prefix
	Children map[string]*Prefix
}

type LabelSearchTerm struct {
	Key   string
	Value *string
}

func (p *Prefix) match(term *LabelSearchTerm) bool {
	for k, v := range p.Labels {
		if k == term.Key && (term.Value == nil || v == *term.Value) {
			return true
		}
	}

	return false
}

// FindPrefixesByLabelTerms fetches every Prefix that has the single 'key' and 'value' label pair set by traversing down the tree starting from whatever value `c` is.
func (p *Prefix) FindPrefixesByLabelTerms(terms []*LabelSearchTerm) []*Prefix {
	parent := p
	result := make([]*Prefix, 0)

	parentMatches := true
	for _, term := range terms {
		if !parent.match(term) {
			parentMatches = false
		}
	}
	if parentMatches {
		result = append(result, parent)
	}

	for _, child := range parent.Children {
		// Note that we now call this method on the actual child of the parent; tree recursion.
		result = append(result, child.FindPrefixesByLabelTerms(terms)...)
	}

	return result
}

func (p *Prefix) PrefixById(id string) *Prefix {
	if p.Cidr == Cidr(id) {
		return p
	}

	for _, alias := range p.Aliases {
		if alias == id {
			return p
		}
	}

	for _, child := range p.Children {
		prefixFromChild := child.PrefixById(id)
		if prefixFromChild != nil {
			return prefixFromChild
		}
	}

	return nil
}

// FindPrefixesByLabelNamesValues fetches every Prefix that has all the 'key' and 'value' label pairs set that are passed through 'l'. Every item of slice 'l' contains a key&value pair following the 'key=value' notation. This Function starts traversing down the tree from whatever value 'c' is.
//func (c *Prefix) FindPrefixesByLabelNamesValues(l []string) []*Prefix {
//	labels :=
//
//	return nil
//}

// AvailableIpSpace returns the available IP space within `c`. C doesn't have to be the parent, if c is a child, this func will return the available IP space in that child.
func (p *Prefix) AvailableIpSpace() ([]Cidr, error) {
	// TODO unfinished
	usedSpace := make([]Cidr, 0, len(p.Children))
	//availableCidrs := []Cidr{}

	for _, child := range p.Children {
		usedSpace = append(usedSpace, child.Cidr)
	}

	sort.Slice(usedSpace, func(i, j int) bool {
		c1 := usedSpace[i]
		c2 := usedSpace[j]

		c1Min, _, err := c1.MinMaxIP()
		if err != nil {
			return false
		}

		c2Min, _, err := c2.MinMaxIP()
		if err != nil {
			return false
		}

		compareCidrs, err := iputil.CompareIPs(c1Min, c2Min)
		if err != nil {
			return false
		}

		return compareCidrs < 0
	})

	minPrefix, maxPrefix, err := p.Cidr.MinMaxIP()
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	minAvailable := minPrefix

	for _, cidr := range usedSpace {
		minUsed, maxUsed, err := cidr.MinMaxIP()
		if err != nil {
			return nil, err
		}

		minAvailableCmpMinUsed, err := iputil.CompareIPs(minAvailable, minUsed)
		if err != nil {
			return nil, err
		}

		if minAvailableCmpMinUsed < 0 {
			fmt.Printf("available: %v - %v \n", minAvailable, minUsed)
			//availableCidrs := append(availableCidrs, cidr)
		}

		// Plus 1 (+1), next available IP
		minAvailable = maxUsed
	}

	minAvailableCmpMaxPrefix, err := iputil.CompareIPs(minAvailable, maxPrefix)
	if err != nil {
		return nil, err
	}

	if minAvailableCmpMaxPrefix < 0 {
		fmt.Printf("available: %v - %v \n", minAvailable, maxPrefix)
		//availableCidrs := append(availableCidrs, cidr)
	}

	return nil, nil
}
