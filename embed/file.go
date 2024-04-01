package embed

import (
	"bytes"
	"io"
	"io/fs"
	"time"

	"go.sancus.dev/core/errors"
)

var (
	_ fs.FileInfo = (*FileInfo)(nil)
	_ fs.DirEntry = (*FileInfo)(nil)
	_ Embeddable  = (*File)(nil)
)

type FileInfo struct {
	sys  *File
	size int64
}

func (fi FileInfo) IsDir() bool       { return false }
func (fi FileInfo) Type() fs.FileMode { return 0 }
func (fi FileInfo) Mode() fs.FileMode { return fs.FileMode(0444) }

func (fi FileInfo) ModTime() time.Time { return fi.sys.ModTime }
func (fi FileInfo) Name() string       { return fi.sys.Name }
func (fi FileInfo) Sys() interface{}   { return fi.sys }

func (fi FileInfo) Size() int64 {
	if fi.size < 0 {
		return fi.sys.Size
	}

	return fi.size
}

func (fi *FileInfo) Info() (fs.FileInfo, error) {
	return fi, nil
}

type Content struct {
	Encoding string
	Bytes    []byte
}

type File struct {
	Name        string
	ModTime     time.Time
	Size        int64
	ContentType string
	Digest      []string
	Content     []Content
}

func (file *File) Info() (fs.FileInfo, error) {
	return FileInfo{file, -1}, nil
}

func (file *File) Encodings() []string {
	var list []string

	for _, c := range file.Content {
		list = append(list, c.Encoding)
	}

	return list
}

func (file *File) NewEncodedReader(encoding string) (io.Reader, int, error) {
	// direct
	for _, c := range file.Content {
		if encoding == c.Encoding {
			// match
			r := bytes.NewReader(c.Bytes)
			return r, len(c.Bytes), nil
		}
	}

	// TODO: indirect

	// not supported
	err := &fs.PathError{
		Path: file.Name,
		Op:   "read",
		Err:  errors.New("Encoding %q not supported", encoding),
	}

	return nil, 0, err
}

func (file *File) Open() (fs.File, error) {
	fd := &FileDescriptor{
		file: file,
		size: -1,
	}
	return fd, nil
}

type FileDescriptor struct {
	encoding string
	reader   io.Reader
	file     *File
	size     int64
}

func (fd *FileDescriptor) Close() error {
	fd.reader = nil
	fd.file = nil
	return nil
}

func (fd *FileDescriptor) Read(buf []byte) (int, error) {
	if fd.reader == nil {
		r, _, err := fd.file.NewEncodedReader(fd.encoding)
		if err != nil {
			return 0, err
		}
		fd.reader = r
	}
	return fd.reader.Read(buf)
}

func (fd *FileDescriptor) Stat() (fs.FileInfo, error) {
	return FileInfo{fd.file, fd.size}, nil
}

func (fd *FileDescriptor) WithEncoding(encoding string) (*FileDescriptor, error) {

	r, l, err := fd.file.NewEncodedReader(encoding)
	if err != nil {
		return nil, err
	}

	fd2 := &FileDescriptor{
		encoding: encoding,
		reader:   r,
		file:     fd.file,
		size:     int64(l),
	}

	return fd2, nil
}
