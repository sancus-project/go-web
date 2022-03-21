package embed

import (
	"io/fs"
)

var (
	_ fs.FS     = (*FS)(nil)
	_ fs.StatFS = (*FS)(nil)
)

type Embeddable interface {
	Info() (fs.FileInfo, error)
	Open() (fs.File, error)
}

type FS struct {
	files map[string]Embeddable
}

func (fsys *FS) Add(name string, file Embeddable) error {
	if file == nil {
		return fs.ErrInvalid
	}

	if fsys.files == nil {
		fsys.files = make(map[string]Embeddable)
	}

	fsys.files[name] = file
	return nil
}

func (fsys *FS) get(name, op string) (Embeddable, *fs.PathError) {
	if f, ok := fsys.files[name]; ok {
		return f, nil
	}

	err := &fs.PathError{
		Op:   op,
		Path: name,
		Err:  fs.ErrNotExist,
	}
	return nil, err
}

func (fsys *FS) Open(name string) (fs.File, error) {
	if f, err := fsys.get(name, "open"); err != nil {
		return nil, err
	} else {
		return f.Open()
	}
}

func (fsys *FS) Stat(name string) (fs.FileInfo, error) {
	if f, err := fsys.get(name, "stat"); err != nil {
		return nil, err
	} else {
		return f.Info()
	}
}
