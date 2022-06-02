package main

import (
	"fmt"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"go.sancus.dev/web/embed"
)

const (
	DefaultPackageName  = "unknown"
	DefaultVariableName = "Files"
	DefaultOutputFile   = "-"
)

type Embedder struct {
	Package string
	Varname string
	Names   []string
	Files   map[string]*embed.File
}

func NewEmbedder(c *Config) (*Embedder, error) {

	m := &Embedder{
		Package: c.Package,
		Varname: c.Varname,
	}

	if err := m.SetDefaults(); err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func ToIdentifier(value, fallback string) (string, bool) {
	if s := os.ExpandEnv(value); s == "" {
		return fallback, true
	} else if token.IsIdentifier(s) {
		return s, true
	} else {
		return fallback, false
	}
}

func (m *Embedder) SetDefaults() error {
	if s, ok := ToIdentifier(m.Package, DefaultPackageName); ok {
		m.Package = s
	} else {
		return fmt.Errorf("invalid %s: %q", "package name", m.Package)
	}

	if s, ok := ToIdentifier(m.Varname, DefaultVariableName); ok {
		m.Varname = s
	} else {
		return fmt.Errorf("invalid %s: %q", "variable name", m.Varname)
	}

	if m.Files == nil {
		m.Files = make(map[string]*embed.File)
	}

	return nil
}

func (m *Embedder) postProcess(name string, blob *Blob) error {
	var s, s0, s1 string

	if s = blob.ContentType; s == "" {
		s0 = "application/octet-stream"
	} else if n := strings.IndexRune(s, ';'); n < 0 {
		s0 = s
	} else {
		s0 = s[:n]
		s1 = s[n+1:]
	}

	switch s0 {
	case "text/plain":
		switch path.Ext(name) {
		case ".css":
			s0 = "text/css"
		case ".js":
			s0 = "application/javascript"
		case ".mjs":
			s0 = "text/javascript"
		case ".hbs":
			s0 = "text/x-handlebars-template"
		default:
			// other text/plain files
			return nil
		}

	default:
		// other files
		return nil
	}

	if s1 != "" {
		s0 += ";" + s1
	}

	if s != s0 {
		blob.ContentType = s0
	}

	return nil
}

func (m *Embedder) addFile(name string, f *os.File, fi fs.FileInfo) error {
	blob := NewBlob()

	if n, err := blob.ReadFrom(f); err != nil {
		return err
	} else if err := m.postProcess(name, blob); err != nil {
		return err
	} else if out, err := blob.Export(fi); err != nil {
		err = &fs.PathError{
			Path: name,
			Op:   "Add",
			Err:  err,
		}
		return err

	} else if n != out.Size {
		err = fmt.Errorf("inconsistent size (%v vs %v)", n, out.Size)
		err = &fs.PathError{
			Path: name,
			Op:   "Add",
			Err:  err,
		}
		return err
	} else {
		name = filepath.ToSlash(name)
		m.Files[name] = out

		log.Printf("%s %7v %s", out.Digest[0], out.Size, name)
		return nil
	}
}

func (m *Embedder) addEntry(base string, name string, fi fs.FileInfo) error {
	// full name
	name = filepath.Join(base, name)

	if fi.Mode().IsRegular() {
		// add file
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()
		return m.addFile(name, f, fi)
	} else if !fi.IsDir() {
		// skip other types
		return nil
	} else if entries, err := os.ReadDir(name); err != nil {
		// readdir error
		return err
	} else {
		return m.addEntries(name, entries)
	}
}

func (m *Embedder) addEntries(base string, entries []fs.DirEntry) error {

	for _, de := range entries {
		if fi, err := de.Info(); err != nil {
			return err
		} else if err := m.addEntry(base, de.Name(), fi); err != nil {
			return err
		}
	}

	return nil
}

func (m *Embedder) Add(name string) error {

	// sanitise
	name = filepath.Clean(name)

	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	if fi, err := f.Stat(); err != nil {
		// not found
		return err
	} else if fi.Mode().IsRegular() {
		// regular
		return m.addFile(name, f, fi)
	} else if fi.IsDir() {
		// directory
		entries, err := f.ReadDir(0)
		if err != nil {
			return err
		}
		return m.addEntries(name, entries)
	} else {
		// invalid type
		return &fs.PathError{
			Path: name,
			Op:   "Add",
			Err:  fs.ErrInvalid,
		}
	}
}
