package kml

import (
	"strings"
	"testing"
)

// TODO negative test cases

func mustParse(t *testing.T, s string) (root *Node, children []*Node) {
	root, err := Parse(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}
	if root == nil {
		t.Fatal("root is nil")
	}

	return root, root.Children
}

func Test_EmptyAndBasics(t *testing.T) {
	root, nodes := mustParse(t, `# comments and the blank lines are ignored

       # comments can be indented

empty:`)

	if len(nodes) != 1 {
		t.Fatalf("invalid number of nodes")
	}

	empty := nodes[0]
	if empty.Key != "empty" {
		t.Fatalf("expected empty, got %s", empty.Key)
	}
	if len(empty.Children) != 0 {
		t.Fatal("empty got children")
	}
	if len(empty.Attributes) != 0 {
		t.Fatal("empty got attributes")
	}
	if empty.Comment != "" {
		t.Fatal("empty got comment")
	}
	if len(empty.Aliases) != 0 {
		t.Fatalf("empty got aliases")
	}
	if empty.Depth != 0 {
		t.Fatal("empty depth != 0")
	}
	if empty.Parent != root {
		t.Fatal("empty parent not root")
	}
	if empty.LineNum != 5 {
		t.Fatal("empty line number not 6")
	}
}

func assertAlias(t *testing.T, node *Node, expected ...string) {
	if len(expected) != len(node.Aliases) {
		t.Fatalf("invalid %s aliases: expected %v, actual %v", node.Key, expected, node.Aliases)
	}

	for index, expectedValue := range expected {
		actualValue := node.Aliases[index]
		if actualValue != expectedValue {
			t.Fatalf("invalid %s aliases: expected %v, actual %v", node.Key, expected, node.Aliases)
		}
	}
}

func Test_Alias(t *testing.T) {
	_, nodes := mustParse(t, `
alias-1: alias
alias-multi-trim:   alias|blias|clias
`)

	assertAlias(t, nodes[0], "alias")
	assertAlias(t, nodes[1], "alias", "blias", "clias")
}

func assertAttr(t *testing.T, node *Node, expectedList ...string) {
	actual := node.Attributes

	expected := make(map[string]string, len(expectedList)/2)
	for i := 0; i < len(expectedList); i += 2 {
		expected[expectedList[i]] = expectedList[i+1]
	}

	if len(actual) != len(expected) {
		t.Fatalf("invalid %s attributes: expected %v, got %v", node.Key, expected, actual)
	}

	for k, expectedValue := range expected {
		actualValue, actualKnown := actual[k]
		if !actualKnown {
			t.Fatalf("invalid %s attributes: expected %v, got %v", node.Key, expected, actual)
		}
		if actualValue != expectedValue {
			t.Fatalf("invalid %s attributes: expected %v, got %v", node.Key, expected, actual)
		}
	}
}

func Test_Attr(t *testing.T) {
	_, nodes := mustParse(t, `
attr-empty: []
attr-1: [key=value]
attr-multi-trim: [ k1 = v1  ,  k2 = v2  ,  k3 = v3 ]
attr-bool: [disabled]
`)
	if len(nodes) != 4 {
		t.Fatal("invalid number of nodes")
	}

	assertAttr(t, nodes[0])
	assertAttr(t, nodes[1], "key", "value")
	assertAttr(t, nodes[2], "k1", "v1", "k2", "v2", "k3", "v3")
	assertAttr(t, nodes[3], "disabled", "true")
}

func assertComment(t *testing.T, node *Node, expected string) {
	if node.Comment != expected {
		t.Fatalf("invalid %s comment: expected '%s', actual '%s'", node.Key, expected, node.Comment)
	}
}

func Test_Comment(t *testing.T) {
	_, nodes := mustParse(t, `
comment: #todo
comment-trim: #   reserved   for future use
`)

	assertComment(t, nodes[0], "todo")
	assertComment(t, nodes[1], "reserved   for future use")
}

func Test_Combos(t *testing.T) {
	_, nodes := mustParse(t, `
# indentation for better readability
alias-attr:         mgmt [deprecated]
alias-comment:      mgmt              #todo remove
alias-attr-comment: mgmt [deprecated] #todo remove
attr-comment: [deprecated] #todo remove
`)

	assertAlias(t, nodes[0], "mgmt")
	assertAttr(t, nodes[0], "deprecated", "true")
	assertComment(t, nodes[0], "")

	assertAlias(t, nodes[1], "mgmt")
	assertAttr(t, nodes[1])
	assertComment(t, nodes[1], "todo remove")

	assertAlias(t, nodes[2], "mgmt")
	assertAttr(t, nodes[2], "deprecated", "true")
	assertComment(t, nodes[2], "todo remove")

	assertAlias(t, nodes[3])
	assertAttr(t, nodes[3], "deprecated", "true")
	assertComment(t, nodes[3], "todo remove")
}

func Test_Indent(t *testing.T) {
	_, nodes := mustParse(t, `
root:
  child1:
    grandchild1:
  child2:
# sneaky comment
    grandchild2:
un-indent-2:
`)

	if len(nodes) != 2 {
		t.Fatal("invalid number of nodes")
	}

	root := nodes[0]
	if len(root.Children) != 2 {
		t.Fatal("invalid number of root children")
	}
	child1 := root.Children[0]
	if len(child1.Children) != 1 && child1.Children[0].Key != "grandchild1" {
		t.Fatal("invalid grandchild1")
	}
	child2 := root.Children[1]
	if len(child2.Children) != 1 && child2.Children[0].Key != "grandchild2" {
		t.Fatal("invalid grandchild2")
	}
}
