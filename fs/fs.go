package fs

import (
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"

	"go.sancus.dev/web"
)

type FileHandler struct {
	fs.File
}

type FS struct {
	FS          fs.FS
	Index       string
	Resolver    ResolverFunc
	ErrNotFound http.Handler
}

func (fsys *FS) getHandler(name, index string) http.Handler {
	if f, _ := fsys.FS.Open(name[1:]); f != nil {
		// match
		if h, ok := f.(http.Handler); ok {
			// implements http.Handler already
			return h
		} else {
			// wrap it to handle
			return &FileHandler{f}
		}
	} else if index != "" {
		// try the index file instead,
		return fsys.getHandler(path.Join(name, index), "")
	} else {
		// Not Found
		return nil
	}
}

func (fsys *FS) handle(rw http.ResponseWriter, req *http.Request, next http.Handler) {
	var name string

	if fsys.Resolver != nil {
		name = fsys.Resolver(req)
	} else {
		name = DefaultResolver(req)
	}

	h := fsys.getHandler(name, fsys.Index)

	if f, ok := h.(io.Closer); ok {
		defer f.Close()
	}

	if h != nil {
		// found
	} else if next != nil {
		// next middleware
		h = next
	} else {
		// standard 404
		h = http.NotFoundHandler()
	}

	h.ServeHTTP(rw, req)
}

var (
	_ web.MiddlewareHandler = (*FS)(nil)
)

func (fsys *FS) Middleware(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		fsys.handle(rw, req, next)
	}

	return http.HandlerFunc(fn)
}

var (
	_ http.Handler = (*FileHandler)(nil)
	_ http.Handler = (*FS)(nil)
)

func (h *FileHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// directly via fs.File?
	if h, ok := h.File.(http.Handler); ok {
		h.ServeHTTP(rw, req)
		return
	}

	// via fs.FileInfo?
	fi, err := h.File.Stat()
	if err != nil {
		log.Panic(err)
	}

	sys := fi.Sys()
	if h, ok := sys.(http.Handler); ok {
		h.ServeHTTP(rw, req)
		return
	}

	// Copy
	if _, err = io.Copy(rw, h.File); err != nil {
		log.Panic(err)
	}
}

func (fsys *FS) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fsys.handle(rw, req, fsys.ErrNotFound)
}
