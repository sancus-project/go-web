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
	sys *File
}

func (fi FileInfo) IsDir() bool       { return false }
func (fi FileInfo) Type() fs.FileMode { return 0 }
func (fi FileInfo) Mode() fs.FileMode { return fs.FileMode(0444) }

func (fi FileInfo) ModTime() time.Time { return fi.sys.ModTime }
func (fi FileInfo) Name() string       { return fi.sys.Name }
func (fi FileInfo) Size() int64        { return fi.sys.Size }
func (fi FileInfo) Sys() interface{}   { return fi.sys }

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
	return FileInfo{file}, nil
}

func (file *File) Encodings() []string {
	var list []string

	for _, c := range file.Content {
		list = append(list, c.Encoding)
	}

	return list
}

func (file *File) NewEncodedReader(encoding string) (io.Reader, error) {
	// direct
	for _, c := range file.Content {
		if encoding == c.Encoding {
			// match
			r := bytes.NewReader(c.Bytes)
			return r, nil
		}
	}

	// TODO: indirect

	// not supported
	err := &fs.PathError{
		Path: file.Name,
		Op:   "read",
		Err:  errors.New("Encoding %q not supported", encoding),
	}

	return nil, err
}

func (file *File) Open() (fs.File, error) {
	fd := &FileDescriptor{
		file: file,
	}
	return fd, nil
}

type FileDescriptor struct {
	encoding string
	reader   io.Reader
	file     *File
}

func (fd *FileDescriptor) Close() error {
	fd.reader = nil
	fd.file = nil
	return nil
}

func (fd *FileDescriptor) Read(buf []byte) (int, error) {
	if fd.reader == nil {
		r, err := fd.file.NewEncodedReader(fd.encoding)
		if err != nil {
			return 0, err
		}
		fd.reader = r
	}
	return fd.reader.Read(buf)
}

func (fd *FileDescriptor) Stat() (fs.FileInfo, error) {
	return fd.file.Info()
}

func (fd *FileDescriptor) WithEncoding(encoding string) (*FileDescriptor, error) {

	r, err := fd.file.NewEncodedReader(encoding)
	if err != nil {
		return nil, err
	}

	fd2 := &FileDescriptor{
		encoding: encoding,
		reader:   r,
		file:     fd.file,
	}
	return fd2, nil
}
