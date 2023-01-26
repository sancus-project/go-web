package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"text/template"
)

func (embed *Embedder) WriteTo(out io.Writer) (int64, error) {
	var buf bytes.Buffer

	// functions
	tmplFuncs := template.FuncMap{
		"formatFile": formatFile,
	}

	// sort filenames
	i := 0
	names := make([]string, len(embed.Files))
	for k, f := range embed.Files {
		if f != nil {
			names[i] = k
			i++
		}
	}

	if i > 0 {
		names = names[:i]
		sort.Strings(names)
	} else {
		names = names[:0]
	}
	embed.Names = names

	// template
	tmpl := `package {{.Package}}
{{- $length := len .Names }}
{{- $varname := .Varname }}
{{- $files := .Files }}

import (
{{- if gt $length 0}}
	"time"

	"go.sancus.dev/web/embed"
{{- else }}
	"go.sancus.dev/web/embed"
{{- end }}
)

var {{ $varname }} embed.FS
{{- if gt $length 0 }}

func init() {
{{- range $name := .Names }}
	{{ $varname }}.Add({{ $name | printf "%q" }},
		{{ index $files $name | formatFile }})
{{- end }}
}
{{- end }}
`

	if t, err := template.New("").Funcs(tmplFuncs).Parse(tmpl); err != nil {
		return 0, err
	} else if err := t.Execute(&buf, embed); err != nil {
		return 0, err
	} else if pretty, err := format.Source(buf.Bytes()); err != nil {
		return 0, err
	} else {
		buf.Reset()
		buf.Write(pretty)

		return io.Copy(out, &buf)
	}
}

func (embed *Embedder) WriteFile(fname string) error {
	// temporary output on the same directory for atomicity
	dirname, basename := filepath.Split(fname)
	if dirname == "" {
		dirname = "."
	}

	if basename == "" {
		return &fs.PathError{
			Path: fname,
			Op:   "open",
			Err:  fs.ErrInvalid,
		}
	}

	// dir/foo.go -> dir/.foo.go~
	tmpname := fmt.Sprintf(".%s~", basename)
	tmpname = filepath.Join(dirname, tmpname)

	f, err := os.OpenFile(tmpname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer os.Remove(tmpname)

	if _, err := embed.WriteTo(f); err != nil {
		// Write error
		return err
	} else if err = f.Sync(); err != nil {
		// Sync error
		return err
	} else if err = f.Close(); err != nil {
		// Close error
		return err
	}

	if err = os.Rename(tmpname, fname); err != nil {
		// Rename error
		return err
	} else {
		// Success
		return nil
	}
}
