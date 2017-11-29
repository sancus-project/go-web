package pathparser

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type nodeType uint

const (
	nodeLiteral nodeType = iota
	nodeCaptureOption
	nodeCaptureIdentifier
	nodeSpecial
)

type Node struct {
	t nodeType
	s string
}

type Path struct {
	nodes []Node
}

func (p *Path) addNode(t nodeType, text string) {
	p.nodes = append(p.nodes, Node{t, text})
}

func (p *Path) beginOptional() {
	p.addNode(nodeSpecial, "[")
}

func (p *Path) beginCapture() {
	p.addNode(nodeSpecial, "{")
}

func (p *Path) endOptional() {
	var np *Node
	var s []string
	var v string
	l := len(p.nodes)
	i := l - 1

	for {
		np = &p.nodes[i]
		if np.t == nodeLiteral {
			s = append(s, np.s)
			i--
		} else {
			break
		}
	}

	v = fmt.Sprintf("(%s)?", strings.Join(s, ""))
	if np = &p.nodes[i-1]; np.t == nodeLiteral && np.s == "/" {
		v = fmt.Sprintf("(/%s)?", v)
		i--
	}
	p.nodes = p.nodes[:i]
	p.addNode(nodeLiteral, v)
}

func (p *Path) endCapture() {
	var s []string
	var v string
	l := len(p.nodes)
	i := l - 1

	for {
		if p.nodes[i].t == nodeCaptureOption {
			s = append(s, p.nodes[i].s)
			i--
		} else {
			break
		}
	}

	switch len(s) {
	case 0:
		v = "[^/]+"
	case 1:
		v = s[0]
	default:
		v = fmt.Sprintf("(%s)", strings.Join(s, "|"))
	}

	v = fmt.Sprintf("(?P<%s>%s)", p.nodes[i].s, v)
	p.nodes = p.nodes[:i-1]
	p.addNode(nodeLiteral, v)
}

func (p *Path) finish() {
	s := []string{`^`}

	for _, np := range p.nodes {
		s = append(s, np.s)
	}

	i, last := 1, len(s)-1
	if s[last] != `/` {
		s = append(s, `$`)
	} else {
		s[last] = `(/|$)`
		i++
	}

	p.nodes[0].s = strings.Join(s, "")
	p.nodes = p.nodes[0:i]
}

func (p *Path) Result() (string, bool) {
	var leaf = (len(p.nodes) == 1)
	return p.nodes[0].s, leaf
}

// Turn path pattern into a regular expression
func MustCompile(path string) (*regexp.Regexp, bool) {
	peg := &Peg{Buffer: path}
	peg.Init()

	if err := peg.Parse(); err != nil {
		log.Fatal(err)
	}

	peg.Execute()
	r, leaf := peg.Result()
	return regexp.MustCompile(r), leaf
}
