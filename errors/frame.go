package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type PC uintptr

func (pc PC) pc() uintptr {
	return uintptr(pc) - 1
}

func (pc PC) Func() *runtime.Func {
	return runtime.FuncForPC(pc.pc())
}

func (pc PC) Name() string {
	if fn := pc.Func(); fn != nil {
		return fn.Name()
	}
	return "unknown"
}

// Heavily based on github.com/pkg/errors.Frame
type Frame struct {
	pc   PC
	file string
	line int
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
			io.WriteString(s, f.pc.Name())
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.file)
		default:
			io.WriteString(s, path.Base(f.file))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(f.line))
	case 'n':
		n := f.pc.Name()
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

func Here(n int) *Frame {
	if pc, file, line, ok := runtime.Caller(n + 1); ok {
		return &Frame{
			pc:   PC(pc),
			file: file,
			line: line,
		}
	}
	return nil
}

func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
