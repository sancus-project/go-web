package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"go.sancus.dev/core/errors"
	"go.sancus.dev/web"
	"go.sancus.dev/web/mimeparse"
	"go.sancus.dev/web/tools"
)

var (
	_ http.Handler = (*ErrorDescriptor)(nil)
	_ web.Handler  = (*ErrorDescriptor)(nil)
	_ web.Error    = (*ErrorDescriptor)(nil)
)

type ErrorDescriptor struct {
	Code    int            `json:"statusCode"`
	Message string         `json:"statusMessage"`
	Header  http.Header    `json:"-"`
	Fatal   error          `json:"fatal,omitempty"`
	Err     []error        `json:"error,omitempty"`
	Stack   []errors.Frame `json:"stack,omitempty"`
}

func (desc *ErrorDescriptor) Status() int {
	return desc.Code
}

func (desc *ErrorDescriptor) Error() string {
	return ErrorText(desc.Code)
}

func (desc *ErrorDescriptor) String() string {
	return ErrorText(desc.Code)
}

// Serve Error as HTTP Response if it's not an actual error
func (desc *ErrorDescriptor) TryServeHTTP(rw http.ResponseWriter, req *http.Request) error {
	switch {
	case desc.Code > 299:
		// error
		return desc
	default:
		// success
		desc.ServeHTTP(rw, req)
		return nil
	}
}

// Serve Error as HTTP Response
func (desc *ErrorDescriptor) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Status
	code := desc.Status()

	// Content-Type
	supported := []string{"text/plain", "application/json"}
	mimetype := mimeparse.BestMatch(supported, req.Header.Get("Accept"))
	if mimetype == "" {
		mimetype = supported[0]
	}

	// Headers
	hdr := rw.Header()
	for k, v := range desc.Header {
		for _, s := range v {
			rw.Header().Add(k, s)
		}
	}

	switch {
	case code == http.StatusOK || code == http.StatusNoContent:
		// quick success
		rw.WriteHeader(http.StatusNoContent)
		return

	default:
		var buf *bytes.Buffer
		var err error

		// error
		tools.SetHeader(hdr, "Content-Type", "%s; charset=utf-8", mimetype)
		tools.SetHeader(hdr, "X-Content-Type-Options", "nosniff")
		rw.WriteHeader(code)

		// content
		switch mimetype {
		case "application/json":
			var b []byte
			b, err = desc.renderJSON()
			if err == nil {
				buf = bytes.NewBuffer(b)
			}
		case "text/plain":
			var b []byte
			buf = bytes.NewBuffer(b)
			err = desc.renderTXT(buf)
		}

		if err != nil {
			// render error
			if buf != nil {
				log.Printf("%+v: %q", errors.Here(), buf.String())
			}
			log.Panic(err)
		} else if buf != nil {
			_, err = io.Copy(rw, buf)
			if err != nil {
				//
				log.Printf("%+v: %s", errors.Here(), err.Error())
			}
		}
	}
}

func (desc *ErrorDescriptor) renderJSON() ([]byte, error) {
	desc.Message = http.StatusText(desc.Code)

	return json.MarshalIndent(desc, "", "  ")
}

func (desc *ErrorDescriptor) renderTXT(w io.Writer) error {
	// Title
	fmt.Fprintln(w, ErrorText(desc.Code))

	// Panic
	if err := desc.Fatal; err != nil {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "panic:", err)
	}

	// Errors
	if len(desc.Err) > 0 {
		fmt.Fprintln(w)

		for i, err := range desc.Err {
			fmt.Fprintf(w, "%v. %s\n", i, err.Error())
		}
	}

	// StackTrace
	if len(desc.Stack) > 0 {
		fmt.Fprintln(w)

		for _, frame := range desc.Stack {
			fmt.Fprintf(w, "%#+v\n", frame)
		}
	}

	return nil
}

func (desc *ErrorDescriptor) AddHeader(key, value string, args ...interface{}) *ErrorDescriptor {
	if len(args) > 0 {
		value = fmt.Sprintf(value, args...)
	}
	value = strings.TrimSpace(value)
	if len(value) > 0 {
		desc.Header.Add(key, value)
	}

	return desc
}

func (desc *ErrorDescriptor) AddErrors(errs ...error) *ErrorDescriptor {
	for _, err := range errs {
		if err != nil {
			desc.Err = append(desc.Err, err)
		}
	}
	return desc
}

func AsDescriptor(err web.Error) *ErrorDescriptor {

	code := err.Status()

	// Summary
	desc := &ErrorDescriptor{
		Code:   code,
		Header: make(map[string][]string),
	}

	// Headers
	if he, ok := err.(interface {
		Headers() http.Header
	}); ok {
		for k, v := range he.Headers() {
			switch k {
			case "Context-Type", "X-Context-Type-Options":
				// skip
			default:
				for _, s := range v {
					desc.AddHeader(k, s)
				}
			}
		}
	}

	if p, ok := err.(interface {
		Recovered() error
	}); ok {
		// Panic
		desc.Fatal = p.Recovered()
	} else if p, ok := err.(interface {
		Errors() []error
	}); ok {
		// Validator
		desc.AddErrors(p.Errors()...)
	} else if p := errors.Unwrap(err); p != nil {
		// Wrapped
		desc.AddErrors(p)
	}

	// StackTrace
	if p, ok := errors.AsStackTracer(err); ok {
		desc.Stack = p.StackTrace()
	}

	return desc
}
