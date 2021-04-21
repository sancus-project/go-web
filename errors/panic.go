package errors

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"go.sancus.dev/web"
)

type PanicError struct {
	stack []string
	rvr   interface{}
}

func (p PanicError) Error() string {
	return fmt.Sprintf("panic: %s", p.rvr)
}

func (_ PanicError) Status() int {
	return http.StatusInternalServerError
}

func (p PanicError) Unwrap() error {
	if err, ok := p.rvr.(error); ok {
		return err
	} else {
		return nil
	}
}

func (p PanicError) Stack() []string {
	return p.stack
}

func (p PanicError) String() string {
	buf := &bytes.Buffer{}

	fmt.Fprintln(buf, "panic:", p.rvr)
	fmt.Fprintln(buf)

	for _, l := range p.stack {
		fmt.Fprintln(buf, l)
	}

	return buf.String()
}

func (p PanicError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := p.Status()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	fmt.Fprintf(w, "%s (Error %v)\n\n", http.StatusText(code), code)
	fmt.Fprintln(w, "panic:", p.rvr)
	fmt.Fprintln(w)

	for _, l := range p.stack {
		fmt.Fprintln(w, l)
	}
}

// backtrace based on github.com/go-chi/middleware/recoverer
func Panic(rvr interface{}) web.Error {
	// process debug stack info
	stack := strings.Split(string(debug.Stack()), "\n")
	lines := []string{}

	// locate panic line, as we may have nested panics
	for i := len(stack) - 1; i > 0; i-- {
		lines = append(lines, stack[i])
		if strings.HasPrefix(stack[i], "panic(0x") {
			lines = lines[0 : len(lines)-2] // remove boilerplate
			break
		}
	}

	// reverse
	for i := len(lines)/2 - 1; i >= 0; i-- {
		opp := len(lines) - 1 - i
		lines[i], lines[opp] = lines[opp], lines[i]
	}

	v := &PanicError{
		stack: lines,
		rvr:   rvr,
	}

	return v
}
