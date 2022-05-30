package router

import (
	"path"
	"regexp"
	"strings"

	"go.sancus.dev/web"
	"go.sancus.dev/web/context"
	"go.sancus.dev/web/errors"
)

func (mux *Mux) findBestNode(path string, args map[string]string) (string, string, *node) {
	var m nodeMatch

	// literal wins, unless shorter
	if s, v, ok := mux.trie.LongestPrefix(path); ok {
		if h, ok := v.(*node); !ok {
			// wtf, how did this get in the trie?
			panic(errors.New("bad node at %q (%T)", s, v))
		} else if s == path {
			// exact match
			return s, "", h
		} else if strings.HasSuffix(h.Pattern, "/*") {
			// match for foo/* patterns

			// `/*` special case
			if s == "/" {
				s = ""
			}

			if l := len(s); path[l] == '/' {
				// good match
				m.Set(s, path[l:], h, nil, nil)
			}
		}
	}

	// re
	for re, h := range mux.pattern {
		if v := re.FindStringSubmatch(path); v != nil {

			s := v[0]
			l := len(s)
			if l > 0 {
				// match

				if s[l-1] == '/' {
					// remove trailing slash from match
					l--
					s = s[:l]
				}

				// test if better than the literal match
				m.Try(s, path[l:], h, re, v)
			}
		}
	}

	return m.Return(args)
}

type nodeMatch struct {
	Path    string
	Extra   string
	Node    *node
	Expr    *regexp.Regexp
	Matches []string
}

func (m *nodeMatch) Set(s0, s1 string, n *node, re *regexp.Regexp, matches []string) {
	*m = nodeMatch{
		Path:    s0,
		Extra:   s1,
		Node:    n,
		Expr:    re,
		Matches: matches,
	}
}

func (m *nodeMatch) Try(s0, s1 string, n *node, re *regexp.Regexp, matches []string) {
	if len(s0) > len(m.Path) {
		m.Set(s0, s1, n, re, matches)
	}
}

func (m *nodeMatch) Return(args map[string]string) (string, string, *node) {
	if len(m.Path) > 0 {

		// return arguments via parameter if requested
		if args != nil {
			if re := m.Expr; re != nil {
				for _, k := range re.SubexpNames() {
					if j := re.SubexpIndex(k); j > 0 {
						args[k] = m.Matches[j]
					}
				}
			}
		}

		return m.Path, m.Extra, m.Node
	}

	// fail
	return "", "", nil
}

// Resolve finds the best handler for a path and returns the corresponding RouteContext
func (m *Mux) Resolve(path string, rctx *context.RoutingContext) (web.Handler, *context.RoutingContext, bool) {
	args := make(map[string]string)
	s0, s1, h := m.findBestNode(path, args)

	if h == nil {
		return nil, nil, false
	}

	if s0 != "/" && s1 == "" && strings.HasSuffix(h.Pattern, "/*") {
		// redirect to the root of the subrouter
		return errors.NewPermanentRedirect("%s/", m.base(path)), rctx, true
	}

	if rctx != nil {
		rctx = rctx.Step(s0, args)
	} else {
		rctx = context.NewRouteContext(s0, s1)

		for k, v := range args {
			rctx.Set(k, v)
		}
	}

	return h, rctx, true
}

func (m *Mux) base(file string) string {
	return path.Base(file)
}
