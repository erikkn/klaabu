package terraform

import (
	"encoding/json"

	"github.com/transferwise/klaabu/klaabu"
)

type module struct {
	Output map[string]output `json:"output"`
}

type output struct {
	Value map[string]prefix `json:"value"`
}

type prefix struct {
	Cidr   string            `json:"cidr"`
	Labels map[string]string `json:"labels,omitempty"`
}

// Generate takes the schema as input (s) and generates JSON data which can be used by Terraform.
func Generate(s *klaabu.Schema) ([]byte, error) {
	aliases := make(map[string]prefix)

	m := module{
		Output: map[string]output{
			"aliases": output{
				Value: aliases,
			},
		},
	}

	populatePrefixes(s.Root, aliases)
	return json.MarshalIndent(m, "", "  ")
}

func populatePrefixes(p *klaabu.Prefix, aliases map[string]prefix) {
	// We don't want to visualize root (transient field).
	if p.Cidr != klaabu.Cidr("0.0.0.0/0") {
		for _, alias := range p.Aliases {
			aliases[alias] = prefix{Cidr: string(p.Cidr), Labels: p.Labels}
		}
	}

	for _, child := range p.Children {
		populatePrefixes(child, aliases)
	}
}
