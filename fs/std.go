package fs

import (
	"io/fs"
)

func Sub(fsys fs.FS, dir string) (fs.FS, error) {
	return fs.Sub(fsys, dir)
}
