package context

import (
	"fmt"
	"io"
	"strings"
)

// Path returns to the full path in a string, including prefix
func (rctx *RoutingContext) Path() string {
	var s string

	prefix := rctx.RoutePrefix
	path := rctx.RoutePath

	if prefix == "/" {
		s = path
	} else if path == "" {
		s = prefix
	} else {
		s = prefix + path
	}

	return s
}

// String returns to the full path in a string, including prefix
func (rctx *RoutingContext) String() string {
	return rctx.Path()
}

// Format formats the RoutingContext according to the fmt.Formatter interface
//
// %s   string corresponding to the full path, including prefix
// %q   quoted strings corresponding to prefix and path, separated by a comma
// %+q  same as %q but also including the pattern
// %v   same as %+q
// %#v  same as %v but also including the parameters
// %# v same as %#v but parameters are indented
func (rctx *RoutingContext) Format(f fmt.State, verb rune) {
	switch verb {
	case 's':
		io.WriteString(f, rctx.Path())
	case 'q':
		fmt.Fprintf(f, "%q, %q", rctx.RoutePrefix, rctx.RoutePath)
		if f.Flag('+') {
			fmt.Fprintf(f, ", %q", rctx.RoutePattern)
		}
	case 'v':
		fmt.Fprintf(f, "%q, %q, %q", rctx.RoutePrefix, rctx.RoutePath, rctx.RoutePattern)

		if f.Flag('#') {
			params := make([]string, len(rctx.RouteParams), 0)

			for k, v := range rctx.RouteParams {
				var s string

				if v, ok := v.(string); ok {
					s = fmt.Sprintf("%s: %q", k, v)
				} else {
					s = fmt.Sprintf("%s: %v", k, v)
				}
				params = append(params, s)
			}

			if len(params) == 0 {
				io.WriteString(f, " {}")
			} else {
				var start, delim, ending string

				if f.Flag(' ') {
					var indent int

					if n, ok := f.Width(); ok && n > 0 {
						indent = n
					}

					tabs := strings.Repeat("\t", indent+1)
					prefix := "\n" + tabs

					start = " {" + prefix
					delim = "," + prefix
					ending = prefix + "}"
				} else {
					start = " {"
					delim = ", "
					ending = "}"
				}

				io.WriteString(f, start)
				io.WriteString(f, strings.Join(params, delim))
				io.WriteString(f, ending)
			}

		}
	}
}
