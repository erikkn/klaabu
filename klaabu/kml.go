package klaabu

import (
	"errors"
	"fmt"
	"github.com/erikkn/klaabu/klaabu/kml"
	"io"
	"os"
)

const (
	cidrs    = "cidrs"
	universe = "0.0.0.0/0"
)

func MarshalKml(node *kml.Node, w io.Writer) error {
	for _, child := range node.Children {
		if child.Key == universe {
			child.Key = cidrs
		}
	}

	err := node.Marshal(w)
	if err != nil {
		return err
	}

	for _, child := range node.Children {
		if child.Key == cidrs {
			child.Key = universe
		}
	}

	return nil
}

// LoadSchemaFromKmlFile parse Schema out of KML file.
func LoadSchemaFromKmlFile(source string) (*Schema, error) {
	node, err := LoadKmlFromFile(source)
	if err != nil {
		return nil, err
	}

	return KmlToSchema(node)
}

func LoadKmlFromFile(source string) (*kml.Node, error) {
	f, err := os.Open(source)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s", err)
	}

	return kml.Parse(f)
}

func KmlToSchema(rootNode *kml.Node) (*Schema, error) {
	topLevelNodes := make(map[string]*kml.Node)
	for _, node := range rootNode.Children {
		if topLevelNodes[node.Key] != nil {
			return nil, fmt.Errorf("duplicate node: '%s' at line %d", node.Key, node.LineNum)
		}

		topLevelNodes[node.Key] = node
	}

	schemaNode := topLevelNodes["schema"]
	if schemaNode == nil {
		return nil, errors.New("'schema' node not found")
	}
	schemaVersion := schemaNode.Attributes["version"]
	if schemaVersion != "v1" {
		return nil, fmt.Errorf("unsupported schema version: %s", schemaVersion)
	}

	rootPrefixNode := topLevelNodes[cidrs]
	if rootPrefixNode == nil {
		return nil, errors.New("'cidrs' node not found")
	}

	err := detectAliasDuplicates(rootPrefixNode, make(map[string]*kml.Node))
	if err != nil {
		return nil, err
	}

	schema := NewSchema(schemaNode.Attributes)

	rootPrefixNode.Key = universe
	rootPrefix, err := prefixFromKmlNode(rootPrefixNode, nil)
	if err != nil {
		return nil, err
	}
	schema.Root = rootPrefix

	return schema, nil
}

func detectAliasDuplicates(node *kml.Node, nodesByAlias map[string]*kml.Node) error {
	for _, alias := range node.Aliases {
		knownNode := nodesByAlias[alias]
		if knownNode != nil {
			return fmt.Errorf("duplicate alias: %s for %s (line %d) and %s (line %d). TIP: use -auxN suffix when provisioning additional VPC CIDRs", alias, knownNode.Key, knownNode.LineNum, node.Key, node.LineNum)
		}
		nodesByAlias[alias] = node
	}

	for _, child := range node.Children {
		err := detectAliasDuplicates(child, nodesByAlias)
		if err != nil {
			return err
		}
	}

	return nil
}

func prefixFromKmlNode(node *kml.Node, parent *Prefix) (*Prefix, error) {
	cidr := Cidr(node.Key)
	err := cidr.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR at line %d: %v", node.LineNum, err)
	}

	prefix := &Prefix{
		Aliases:  node.Aliases,
		Cidr:     cidr,
		Labels:   node.Attributes,
		Parent:   parent,
		Children: make(map[string]*Prefix, len(node.Children)),
	}

	for _, childNode := range node.Children {
		if prefix.Children[childNode.Key] != nil {
			return nil, fmt.Errorf("duplicate prefix/CIDR: %s at line %d", childNode.Key, childNode.LineNum)
		}

		childPrefix, err := prefixFromKmlNode(childNode, prefix)
		if err != nil {
			return nil, err
		}
		prefix.Children[childNode.Key] = childPrefix
	}

	return prefix, nil
}
