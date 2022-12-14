package kml

import (
	"io"
	"sort"
	"strings"

	"github.com/transferwise/klaabu/klaabu/iputil"
)

type ColumnWidths struct {
	key, aliases, attributes int
}

func (n *Node) Marshal(w io.Writer) error {
	widths := &ColumnWidths{}
	n.calculateColumnWidths(widths)
	return n.marshal(w, widths)
}

func (n *Node) marshal(w io.Writer, widths *ColumnWidths) error {
	if n.Key != "" {
		err := n.marshalNode(w, widths)
		if err != nil {
			return err
		}
	}

	sort.Slice(n.Children, func(i1, i2 int) bool {
		min1, _, _ := iputil.MinMaxIP(n.Children[i1].Key)
		min2, _, _ := iputil.MinMaxIP(n.Children[i2].Key)
		cmp, _ := iputil.CompareIPs(min1, min2)
		return cmp < 0
	})

	for index, child := range n.Children {
		if child.Depth == 0 && index != 0 {
			_, err := w.Write([]byte("\n"))
			if err != nil {
				return err
			}
		}

		err := child.marshal(w, widths)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) marshalNode(w io.Writer, widths *ColumnWidths) error {
	// use builder and avoid intermediate strings to be GC/CPU efficient
	b := strings.Builder{}
	b.Grow(widths.key + widths.aliases + widths.attributes + len(n.Comment))

	const singleIndent = "  "
	for i := 0; i < n.Depth; i++ {
		b.WriteString(singleIndent)
	}
	b.WriteString(n.Key)
	b.WriteString(": ")
	for b.Len() < widths.key {
		b.WriteString(" ")
	}

	const aliasSeparator = "|"
	for i, alias := range n.Aliases {
		if i != 0 {
			b.WriteString(aliasSeparator)
		}
		b.WriteString(alias)
	}
	for b.Len() < widths.key+widths.aliases {
		b.WriteString(" ")
	}

	if len(n.Attributes) > 0 {
		b.WriteString("[")
		attributesWritten := 0
		const attrPairSeparator = ","
		const attrKeyValueSeparator = "="

		// make ordering consistent
		attributeKeys := make([]string, 0, len(n.Attributes))
		for k, _ := range n.Attributes {
			attributeKeys = append(attributeKeys, k)
		}
		sort.Strings(attributeKeys)

		for _, k := range attributeKeys {
			v := n.Attributes[k]
			if attributesWritten > 0 {
				b.WriteString(attrPairSeparator)
			}
			b.WriteString(k)
			b.WriteString(attrKeyValueSeparator)
			b.WriteString(v)

			attributesWritten++
		}
		b.WriteString("]")
	}
	for b.Len() < widths.key+widths.aliases+widths.attributes {
		b.WriteString(" ")
	}

	if len(n.Comment) > 0 {
		b.WriteString(" # ")
		b.WriteString(n.Comment)
	}

	b.WriteString("\n")

	_, err := w.Write([]byte(b.String()))
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) calculateColumnWidths(widths *ColumnWidths) {
	nodeKeyWidth := n.Depth*2 + len(n.Key) + 2 // append ": "
	if nodeKeyWidth > widths.key {
		widths.key = nodeKeyWidth
	}

	nodeAliasesWidth := 0
	for _, alias := range n.Aliases {
		if nodeAliasesWidth > 0 {
			nodeAliasesWidth += 1 // separator
		}
		nodeAliasesWidth += len(alias)
	}
	if nodeAliasesWidth > 0 {
		nodeAliasesWidth += 1 // padding
	}
	if nodeAliasesWidth > widths.aliases {
		widths.aliases = nodeAliasesWidth
	}

	nodeAttrsWidth := 2
	for k, v := range n.Attributes {
		if nodeAttrsWidth > 2 {
			nodeAttrsWidth += 1 // separator
		}
		nodeAttrsWidth += len(k) + 1 + len(v)
	}
	if nodeAttrsWidth > widths.attributes {
		widths.attributes = nodeAttrsWidth
	}

	for _, child := range n.Children {
		child.calculateColumnWidths(widths)
	}
}
