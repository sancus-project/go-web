package intercept

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/felixge/httpsnoop"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

type WriteInterceptor struct {
	code           int
	headersWritten bool

	rw      http.ResponseWriter
	headers http.Header
}

func (m *WriteInterceptor) Writer() http.ResponseWriter {
	return m.rw
}

func (m *WriteInterceptor) Error() web.Error {
	return nil
}

func (m *WriteInterceptor) Header(original httpsnoop.HeaderFunc) http.Header {
	log.Printf("%+n()", errors.Here(0))
	return original()
}

func (m *WriteInterceptor) Write(original httpsnoop.WriteFunc, b []byte) (int, error) {
	log.Printf("%+n(%q (%v))", errors.Here(0), b, len(b))
	return original(b)
}

func (m *WriteInterceptor) WriteHeader(original httpsnoop.WriteHeaderFunc, code int) {
	log.Printf("%+n(%v)", errors.Here(0), code)
	original(code)
}

func (m *WriteInterceptor) Flush(original httpsnoop.FlushFunc) {
	log.Printf("%+n()", errors.Here(0))
	original()
}

func (m *WriteInterceptor) CloseNotify(original httpsnoop.CloseNotifyFunc) <-chan bool {
	log.Printf("%+n()", errors.Here(0))
	return original()
}

func (m *WriteInterceptor) Hijack(original httpsnoop.HijackFunc) (net.Conn, *bufio.ReadWriter, error) {
	log.Printf("%+n()", errors.Here(0))
	return original()
}

func (m *WriteInterceptor) ReadFrom(original httpsnoop.ReadFromFunc, src io.Reader) (int64, error) {
	log.Printf("%+n(%v)", errors.Here(0), src)
	return original(src)
}

func (m *WriteInterceptor) Push(original httpsnoop.PushFunc, target string, opts *http.PushOptions) error {
	log.Printf("%+n(%q, %#v)", errors.Here(0), target, opts)
	return original(target, opts)
}

func NewWriter(w http.ResponseWriter) *WriteInterceptor {

	m := &WriteInterceptor{
		headers: w.Header(),
	}

	hooks := httpsnoop.Hooks{
		Header: func(original httpsnoop.HeaderFunc) httpsnoop.HeaderFunc {
			return func() http.Header {
				return m.Header(original)
			}
		},

		WriteHeader: func(original httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(code int) {
				m.WriteHeader(original, code)
			}
		},

		Write: func(original httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return func(b []byte) (int, error) {
				return m.Write(original, b)
			}
		},

		Flush: func(original httpsnoop.FlushFunc) httpsnoop.FlushFunc {
			return func() {
				m.Flush(original)
			}
		},

		CloseNotify: func(original httpsnoop.CloseNotifyFunc) httpsnoop.CloseNotifyFunc {
			return func() <-chan bool {
				return m.CloseNotify(original)
			}
		},

		Hijack: func(original httpsnoop.HijackFunc) httpsnoop.HijackFunc {
			return func() (net.Conn, *bufio.ReadWriter, error) {
				return m.Hijack(original)
			}
		},

		ReadFrom: func(original httpsnoop.ReadFromFunc) httpsnoop.ReadFromFunc {
			return func(src io.Reader) (int64, error) {
				return m.ReadFrom(original, src)
			}
		},

		Push: func(original httpsnoop.PushFunc) httpsnoop.PushFunc {
			return func(target string, opts *http.PushOptions) error {
				return m.Push(original, target, opts)
			}
		},
	}

	m.rw = httpsnoop.Wrap(w, hooks)
	return m
}
