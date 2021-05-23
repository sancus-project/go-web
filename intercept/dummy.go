package intercept

import (
	"bytes"
	"io"
	"net/http"

	"go.sancus.dev/web/errors"
)

type DummyWriter struct {
	buffer bytes.Buffer
	header http.Header
	code   int

	headersWritten bool
}

func (dw *DummyWriter) Status() int {
	if dw.code == 0 {
		return http.StatusOK
	}
	return dw.code
}

func (dw *DummyWriter) Header() http.Header {
	if dw.header == nil {
		dw.header = make(http.Header)
	}
	return dw.header
}

func (dw *DummyWriter) Write(b []byte) (int, error) {
	if !dw.headersWritten {
		dw.WriteHeader(http.StatusOK)
	}

	return dw.buffer.Write(b)
}

func (dw *DummyWriter) WriteHeader(code int) {
	dw.code = code
	dw.headersWritten = true
}

func (dw *DummyWriter) Error() error {
	return errors.NewError(dw.Status(), dw.Header(), dw.buffer.Bytes())
}

func (dw *DummyWriter) TryServeHTTP(w http.ResponseWriter, r *http.Request) error {
	if err := dw.Error(); err != nil {
		return err
	}

	h := w.Header()
	for k, x := range dw.Header() {
		for _, v := range x {
			h.Add(k, v)
		}
	}
	w.WriteHeader(dw.Status())

	b := bytes.NewBuffer(dw.buffer.Bytes())

	if rrf, ok := w.(io.ReaderFrom); ok {
		_, err := rrf.ReadFrom(b)
		return err
	}

	_, err := b.WriteTo(w)
	return err
}
