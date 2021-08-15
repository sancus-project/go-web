package router

import (
	"regexp"
	"strings"

	"go.sancus.dev/web/pathparser"
)

type parser struct {
	path    string
	pattern string
	re      string
}

func (p *parser) Path() string {
	return p.path
}

func (p *parser) Pattern() string {
	return p.pattern
}

func (p *parser) Literal() bool {
	return len(p.re) == 0
}

func (p *parser) Compile() (*regexp.Regexp, error) {
	if len(p.re) > 0 {
		return regexp.Compile(p.re)
	}
	return nil, nil
}

func (mux *Mux) parsePath(path string) (*parser, error) {

	p := &parser{
		pattern: path,
		path: path,
	}

	if path == "/" {
		// special case, root only
	} else if path == "/*" {
		// special case, root or anything bellow
		p.path = "/"
	} else {
		if strings.HasSuffix(path, "/*") {
			// remove trailing * from foo/* on paths
			p.path = path[0 : len(path)-1]
		} else if strings.HasSuffix(path, "/") {
			// add trailing * to foo/ on pattern
			p.pattern += "*"
		}

		peg := pathparser.Peg{
			Buffer: p.path,
		}
		peg.Init()

		if err := peg.Parse(); err != nil {
			return nil, err
		}
		peg.Execute()

		if !peg.Literal() {
			p.re, _ = peg.Result()
		} else if s := p.path; strings.HasSuffix(s, "/") {
			// for literals we don't want the trailing
			// slash so /foo/ also catches /foo like when
			// using regexps
			p.path = s[0 : len(s)-1]
		}
	}

	return p, nil
}
