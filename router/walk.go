package router

import (
	"strings"
)

// WalkFn is used when walking the tree. Takes the pattern and handler,
// returning if iteration should be terminated.
type WalkFn func(pattern string, h Handler) bool

// Walk calls a given function for each node on the tree
func (mux *Mux) Walk(fn WalkFn) {
	if fn != nil {
		mux.walk("", fn)
	}
}

func walk(prefix string, n *node, fn WalkFn) bool {
	h := n.Handler
	pattern := prefix + n.Pattern

	if h, ok := h.(interface {
		walk(string, WalkFn) bool
	}); ok {
		return h.walk(pattern, fn)
	}

	return fn(pattern, h)

}

func (mux *Mux) walk(prefix string, fn WalkFn) bool {
	var done bool

	prefix = strings.TrimSuffix(prefix, "/*")

	// trie
	mux.trie.Walk(func(path string, v interface{}) bool {
		done = walk(prefix, v.(*node), fn)
		return done
	})

	if !done {
		// re
		for _, n := range mux.pattern {
			if walk(prefix, n, fn) {
				return true // done
			}
		}
	}

	return false // continue
}

func (n *rawNode) walk(prefix string, fn WalkFn) bool {
	h := n.h

	if h, ok := h.(interface {
		walk(string, WalkFn) bool
	}); ok {
		return h.walk(prefix, fn)
	}

	return fn(prefix, h.(Handler))
}
