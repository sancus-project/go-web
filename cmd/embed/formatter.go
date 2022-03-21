package main

import (
	"bytes"
	"io"
	"reflect"

	"go.sancus.dev/core/fmt"
	"go.sancus.dev/web/embed"
)

type printer struct {
	io.Writer

	indent string
}

// writer
func (p *printer) WriteByte(b byte) (int, error) {
	return p.Write([]byte{b})
}

func (p *printer) WriteString(s string) (int, error) {
	return p.Write([]byte(s))
}

func (p *printer) Printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(p, format, args...)
}

// indentation
func (p *printer) in() {
	p.indent += "\t"
}

func (p *printer) out() {
	if l := len(p.indent); l > 1 {
		p.indent = p.indent[:l-1]
	} else {
		p.indent = ""
	}
}

func (p *printer) printIndent() {
	p.Write([]byte(p.indent))
}

// rendering
func (p *printer) printInline(v reflect.Value, showtype bool, str string, x interface{}) {
	if showtype {
		typ := v.Type().String()
		p.Printf("%s("+str+")", typ, x)
	} else {
		p.Printf(str, x)
	}
}

func (p printer) PrintValue(v reflect.Value, showtype bool) {

	// try fmt.GoStringer first
	if v.IsValid() && v.CanInterface() {
		x := v.Interface()
		if goStringer, ok := x.(fmt.GoStringer); ok {
			p.WriteString(goStringer.GoString())
			return
		}
	}

	switch v.Kind() {
	case reflect.String:
		p.printInline(v, showtype, "%q", v.String())
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.printInline(v, showtype, "%d", v.Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		p.printInline(v, showtype, "%d", v.Uint())

	case reflect.Ptr:
		if v.IsNil() {
			p.WriteString("nil")
		} else {
			// "&T..."
			p.WriteByte('&')
			p.PrintValue(v.Elem(), true)
		}

	case reflect.Slice, reflect.Array:
		typ := v.Type().String()

		// "[]T{...}"
		p.WriteString(typ)

		if l := v.Len(); l == 0 {
			// empty
			p.WriteString("{}")
		} else {

			p.WriteString("{\n")
			p.in()

			if typ == "[]uint8" {
				// bytes
				i, j := 0, 0
				for i < l {

					if j == 0 {
						p.printIndent()
					} else {
						p.WriteString(", ")
					}

					p.Printf("0x%02x", v.Index(i).Uint())

					i, j = i+1, j+1
					if j == 8 || i == l {
						p.WriteString(",\n")
						j = 0
					}
				}
			} else {
				// others
				for i := 0; i < l; i++ {
					p.printIndent()
					p.PrintValue(v.Index(i), false)
					p.WriteString(",\n")
				}
			}

			p.out()
			p.printIndent()
			p.WriteString("}")
		}

	case reflect.Struct:
		t := v.Type()
		fields := reflect.VisibleFields(t)
		first := true

		// "T{...}"
		if showtype {
			p.WriteString(t.String())
		}

		p.in()
		for _, f := range fields {
			fval := v.FieldByIndex(f.Index)

			// each non-empty visible field
			if fval.IsZero() {
				continue
			}

			if first {
				// open
				p.WriteString("{\n")
				first = false
			}

			p.printIndent()
			p.Printf("%s: ", f.Name)
			p.PrintValue(fval, false)
			p.WriteString(",\n")
		}
		p.out()

		if first {
			// empty
			p.WriteString("{}")
		} else {
			// close
			p.printIndent()
			p.WriteString("}")
		}

	default:
		panic(v.Kind())
	}

}

func formatFile(file *embed.File) string {
	buf := &bytes.Buffer{}

	// generate
	p := &printer{buf, ""}
	p.PrintValue(reflect.ValueOf(file), true)

	return buf.String()
}
