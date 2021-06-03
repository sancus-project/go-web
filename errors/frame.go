package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type StackTrace []Frame

// Heavily based on github.com/pkg/errors.Frame
type Frame struct {
	pc    uintptr
	entry uintptr
	name  string
	file  string
	line  int
}

func frameForPC(pc uintptr) Frame {
	var entry uintptr
	var name string
	var file string
	var line int

	if fp := runtime.FuncForPC(pc - 1); fp != nil {
		entry = fp.Entry()
		name = fp.Name()
		file, line = fp.FileLine(pc)
	} else {
		name = "unknown"
		file = "unknown"
	}

	return Frame{
		pc:    pc,
		entry: entry,
		name:  name,
		file:  file,
		line:  line,
	}
}

func (f Frame) Name() string {
	return f.name
}

// Format formats the frame according to the fmt.Formatter interface.
//
//    %s    source file
//    %d    source line
//    %n    function name
//    %v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+s   function name and path of source file relative to the compile time
//          GOPATH separated by \n\t (<funcname>\n\t<path>)
//    %+n   full package name followed by function name
//    %+v   equivalent to %+s:%d
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.name)
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.file)
		default:
			io.WriteString(s, path.Base(f.file))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(f.line))
	case 'n':
		n := f.name
		switch {
		case s.Flag('+'):
			io.WriteString(s, n)
		default:
			io.WriteString(s, funcname(n))
		}
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

func (st StackTrace) Format(s fmt.State, verb rune) {

	if s.Flag('#') {
		l := len(st)
		for i, f := range st {
			fmt.Fprintf(s, "\n[%v/%v] ", i, l)
			f.Format(s, verb)
		}
	} else {
		for _, f := range st {
			io.WriteString(s, "\n")
			f.Format(s, verb)
		}
	}
}

func Here() *Frame {
	const depth = 1
	var pcs [depth]uintptr

	if n := runtime.Callers(2, pcs[:]); n > 0 {
		f := frameForPC(pcs[0])
		return &f
	}

	return nil
}

func StackFrame(skip int) *Frame {
	const depth = 32
	var pcs [depth]uintptr

	if n := runtime.Callers(2, pcs[:]); n > skip {
		f := frameForPC(pcs[skip])
		return &f
	}

	return nil
}

func BackTrace(skip int) StackTrace {
	const depth = 32
	var pcs [depth]uintptr
	var st StackTrace

	if n := runtime.Callers(2, pcs[:]); n > skip {
		var frames []Frame

		n = n - skip
		frames = make([]Frame, 0, n)

		for _, pc := range pcs[skip:n] {
			frames = append(frames, frameForPC(pc))
		}

		st = StackTrace(frames)
	}

	return st
}

func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
