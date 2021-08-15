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
		pattern: path, // original
		path:    path, // literal
	}

	// "/" and "/*" are special cases, for all other foo/*
	// is the same as foo/

	if path != "/" {

		if s := p.path; s == "/*" {
			// keep
		} else if strings.HasSuffix(s, "/*") {
			// remove trailing *
			p.path = s[0 : len(s)-1]
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
		}

	}

	return p, nil
}
