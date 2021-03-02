package kml

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	ignorableLineRx = regexp.MustCompile(`^\s*(?:\#.*)?$`)
	nodeLineRx      = regexp.MustCompile(`^(?P<indent>\s*)(?P<key>\S+):\s*(?P<aliases>[^\[#\s]+)?\s*(?:\[(?P<attrs>.*)\])?\s*\#?\s*(?P<comment>.*)?$`)
	sub             map[string]int
)

func init() {
	names := nodeLineRx.SubexpNames()
	sub = make(map[string]int, len(names)-1)
	for i, name := range names {
		sub[name] = i
	}
}

type Node struct {
	Aliases    []string
	Attributes map[string]string
	Children   []*Node
	Comment    string
	Depth      int
	Key        string
	LineNum    int
	Parent     *Node
}

func newNode() *Node {
	return &Node{
		Aliases:    make([]string, 0, 1),
		Attributes: make(map[string]string),
		Children:   make([]*Node, 0, 4),
	}
}

func Parse(reader io.Reader) (*Node, error) {
	root := newNode()
	root.Depth = -1

	parent := root

	s := bufio.NewScanner(reader)
	lineNum := 0

	commonIndentChar := '?'
	commonSingleIndentSize := -1

	for s.Scan() {
		lineNum++

		line := s.Text()
		if ignorableLineRx.MatchString(line) {
			continue
		}

		node := newNode()

		groups := nodeLineRx.FindStringSubmatch(line)
		if groups == nil {
			return nil, fmt.Errorf("invalid kml: line %d: '%s'", lineNum, line)
		}

		// indentation
		indent := groups[sub["indent"]]
		indentSize := 0
		indentChar := '?'
		for _, c := range indent {
			indentSize++
			if indentChar == '?' {
				indentChar = c
			} else if indentChar != c {
				return nil, fmt.Errorf("invalid kml: line %d: inconsistent indent character '%c' %U %q", lineNum, c, c, c)
			}
		}
		if indentChar != '?' {
			if commonIndentChar == '?' {
				commonIndentChar = indentChar
			} else if commonIndentChar != indentChar {
				return nil, fmt.Errorf("invalid kml: line %d: inconsistent indent character '%c' %U %q", lineNum, indentChar, indentChar, indentChar)
			}
			if commonSingleIndentSize == -1 {
				commonSingleIndentSize = indentSize
			}
		}
		if indentSize%commonSingleIndentSize != 0 {
			return nil, fmt.Errorf("invalid kml: line %d: inconsistent indentation size, %d is not multiple of %d", lineNum, indentSize, commonSingleIndentSize)
		}
		node.Depth = indentSize / commonSingleIndentSize
		for node.Depth <= parent.Depth {
			// not checking for parent.Parent == nil as it should not be possible
			parent = parent.Parent
		}
		if node.Depth == parent.Depth+2 {
			if len(parent.Children) == 0 {
				return nil, fmt.Errorf("invalid kml: line %d: invalid indent level", lineNum)
			}
			parent = parent.Children[len(parent.Children)-1]
		}
		if node.Depth != parent.Depth+1 {
			return nil, fmt.Errorf("invalid kml: line %d: invalid indent level", lineNum)
		}

		for _, v := range strings.Split(groups[sub["aliases"]], "|") {
			if v != "" {
				node.Aliases = append(node.Aliases, v)
			}
		}

		attrKvStrings := strings.Split(groups[sub["attrs"]], ",")
		node.Attributes = make(map[string]string, len(attrKvStrings))
		for _, kvString := range attrKvStrings {
			kv := strings.Split(kvString, "=")
			for i, v := range kv {
				kv[i] = strings.TrimSpace(v)
			}

			if len(kv) == 2 {
				node.Attributes[kv[0]] = kv[1]
			} else if len(kv) == 1 && kv[0] != "" {
				node.Attributes[kv[0]] = "true"
			}
		}

		node.Comment = strings.TrimSpace(groups[sub["comment"]])

		node.Key = groups[sub["key"]]
		node.LineNum = lineNum

		node.Parent = parent
		parent.Children = append(parent.Children, node)
	}

	if err := s.Err(); err != nil && err != io.EOF {
		return nil, err
	}

	return root, nil
}
