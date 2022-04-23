package embed

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.sancus.dev/web"
	"go.sancus.dev/web/errors"
	"go.sancus.dev/web/qlist"
)

// http.Handler
var (
	_ http.Handler = (*File)(nil)
	_ http.Handler = (*FileDescriptor)(nil)
)

func (f *File) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fd := &FileDescriptor{
		file: f,
	}
	fd.ServeHTTP(rw, req)
}

func (fd *FileDescriptor) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if err := fd.TryServeHTTP(rw, req); err != nil {
		errors.HandleError(rw, req, err)
	}
}

// web.Handler
var (
	_ web.Handler = (*File)(nil)
	_ web.Handler = (*FileDescriptor)(nil)
)

func (f *File) TryServeHTTP(rw http.ResponseWriter, req *http.Request) error {
	fd := &FileDescriptor{
		file: f,
	}
	return fd.TryServeHTTP(rw, req)
}

func (fd *FileDescriptor) TryServeHTTP(rw http.ResponseWriter, req *http.Request) (err error) {
	switch req.Method {
	case "GET", "HEAD":
		var f io.Reader

		// Headers
		setNotZeroTimeHeader(rw, "Last-Modified", fd.file.ModTime)
		setNotZeroHeader(rw, "Content-Type", fd.file.ContentType)

		if s := strings.Join(fd.file.Digest, ","); s != "" {
			rw.Header().Set("Digest", s)

			s = fd.file.Digest[0]
			if i := strings.IndexRune(s, '='); i >= 0 {
				s = s[i+1:]
			}

			rw.Header().Set("Etag", s)
		}

		if req.Method == "GET" {
			var enc string

			// encoding
			q, err := qlist.ParseQualityHeader(req.Header, "Accept-Encoding")
			if err != nil {
				err = errors.Wrap(err, "Accept-Encoding")
				return errors.BadRequest(err)
			}

			enc, f, err = fd.WithAcceptEncodingQuality(q)
			if err != nil {
				return err
			}

			setNotZeroHeader(rw, "Content-Encoding", enc)
		}

		// only emit Content-Length after the request has been accepted
		setFormattedHeader(rw, "Content-Length", "%d", fd.file.Size)

		if f != nil {
			// GET, copy encoded data
			_, err = io.Copy(rw, f)
		}
		return

	default:
		// OPTIONS
		return errors.MethodNotAllowed(req.Method, "GET", "HEAD")
	}

}

func (fd *FileDescriptor) WithAcceptEncodingQuality(ql qlist.QualityList) (string, *FileDescriptor, error) {
	// choose encoding
	encodings := fd.file.Encodings()
	best, ok := qlist.BestEncodingQuality(encodings, ql)
	if !ok {
		return "", nil, errors.ErrNotAcceptable
	}

	if best == "identity" {
		best = ""
	}

	fd2, err := fd.WithEncoding(best)
	return best, fd2, err
}

func setNotZeroTimeHeader(rw http.ResponseWriter, name string, value time.Time) {
	if value.IsZero() {
		// skip
	} else if value.Location() == nil || value.Location() == time.UTC {
		// UTC
		s := value.Format("Mon, 02 Jan 2006 15:04:05")
		setFormattedHeader(rw, name, "%s GMT", s)
	} else {
		// non-UTC
		s := value.Format(time.RFC1123)
		setHeader(rw, name, s)
	}
}

func setNotZeroHeader(rw http.ResponseWriter, name string, value string) {
	if value != "" {
		setHeader(rw, name, value)
	}
}

func setFormattedHeader(rw http.ResponseWriter, name, format string, args ...interface{}) {
	setHeader(rw, name, fmt.Sprintf(format, args...))
}

func setHeader(rw http.ResponseWriter, name string, value string) {
	rw.Header().Set(name, value)
}
