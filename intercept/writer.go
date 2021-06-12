package intercept

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/felixge/httpsnoop"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
)

var (
	ErrNotImplemented = errors.New("Not Implemented")
)

type WriteInterceptor struct {
	buffer         bytes.Buffer
	code           int
	mute           bool
	capture        bool
	headersWritten bool

	rw       http.ResponseWriter // ResponseWriter wrapper
	header   http.Header         // Working copy of Headers
	original http.Header         // Original Headers table
}

func (m *WriteInterceptor) Writer() http.ResponseWriter {
	return m.rw
}

func (m *WriteInterceptor) Error() web.Error {
	if !m.headersWritten {
		return &errors.HandlerError{
			Code:   http.StatusNoContent,
			Header: m.header,
		}
	}

	return errors.NewWebError(m.code, m.header, m.buffer.Bytes())
}

func (m *WriteInterceptor) Header(original httpsnoop.HeaderFunc) http.Header {
	return m.header
}

func (m *WriteInterceptor) Write(original httpsnoop.WriteFunc, b []byte) (int, error) {
	if !m.headersWritten {
		m.rw.WriteHeader(http.StatusOK)
	}

	if m.capture {
		// buffer
		return m.buffer.Write(b)
	} else if m.mute {
		// fake
		return len(b), nil
	} else {
		// real
		return original(b)
	}
}

func (m *WriteInterceptor) WriteHeader(original httpsnoop.WriteHeaderFunc, code int) {
	if m.headersWritten {
		log.Fatal(errors.New("%+n(%v): %s", errors.Here(), code, "Invalid Call"))
	}

	m.headersWritten = true
	m.code = code

	if code >= http.StatusContinue && code < http.StatusMultipleChoices {
		// good, copy headers and write them

		if code == http.StatusNoContent {
			m.mute = true
		}

		for k, _ := range m.original {
			if w, ok := m.header[k]; !ok {
				// delete deleted headers
				m.original.Del(k)
			} else {
				// replace value of those that remain
				m.original[k] = w
			}
		}

		for k, v := range m.header {
			if _, ok := m.original[k]; !ok {
				// add new headers
				m.original[k] = v
			}
		}

		original(code)

	} else {
		// capture writes for later review
		m.capture = true
	}
}

func (m *WriteInterceptor) Flush(original httpsnoop.FlushFunc) {
	err := ErrNotImplemented
	log.Fatal(err)
}

func (m *WriteInterceptor) CloseNotify(original httpsnoop.CloseNotifyFunc) <-chan bool {
	err := ErrNotImplemented
	log.Fatal(err)
	return nil
}

func (m *WriteInterceptor) Hijack(original httpsnoop.HijackFunc) (net.Conn, *bufio.ReadWriter, error) {
	err := ErrNotImplemented
	log.Fatal(err)
	return nil, nil, err
}

func (m *WriteInterceptor) ReadFrom(original httpsnoop.ReadFromFunc, src io.Reader) (int64, error) {
	err := ErrNotImplemented
	log.Fatal(err)
	return 0, err
}

func (m *WriteInterceptor) Push(original httpsnoop.PushFunc, target string, opts *http.PushOptions) error {
	err := ErrNotImplemented
	log.Fatal(err)
	return err
}

func NewWriter(w http.ResponseWriter, method string) *WriteInterceptor {

	var mute bool

	if method == "HEAD" || method == "OPTIONS" {
		mute = true
	}

	h := w.Header()
	m := &WriteInterceptor{
		original: h,
		header:   h.Clone(),
		mute:     mute,
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
