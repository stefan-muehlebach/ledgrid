//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"
	"path"
	"runtime"
	"text/template"
)







const embedTemplate = `
// {{ .Name }} supports binding a {{ .Type }} value.
type {{ .Name }} interface {
    DataItem
    Get() ({{ .Type }})
    Set({{ .Type }})
}

// External{{ .Name }} supports binding a {{ .Type }} value to an external value.
type External{{ .Name }} interface {
    {{ .Name }}
    Reload()
}

// New{{ .Name }} returns a bindable {{ .Type }} value that is managed internally.
func New{{ .Name }}() {{ .Name }} {
    var blank {{ .Type }} = {{ .Default }}
    b := &bound{{ .Name }}{val: &blank}
    b.Init(b)
    return b
}

// Bind{{ .Name }} returns a new bindable value that controls the contents of the provided {{ .Type }} variable.
// If your code changes the content of the variable this refers to you should call Reload() to inform the bindings.
func Bind{{ .Name }}(v *{{ .Type }}) External{{ .Name }} {
    if v == nil {
        var blank {{ .Type }} = {{ .Default }}
        v = &blank // never allow a nil value pointer
    }
    b := &boundExternal{{ .Name }}{}
    b.val = v
    b.old = *v
    b.Init(b)
    return b
}

type bound{{ .Name }} struct {
    base
    val *{{ .Type }}
}

func (b *bound{{ .Name }}) Get() ({{ .Type }}) {
    b.lock.RLock()
    defer b.lock.RUnlock()
    if b.val == nil {
        return {{ .Default }}
    }
    return *b.val
}

func (b *bound{{ .Name }}) Set(val {{ .Type }}) {
    b.lock.Lock()
    defer b.lock.Unlock()
    {{- if eq .Comparator "" }}
    if *b.val == val {
        return
    }
    {{- else }}
    if {{ .Comparator }}(*b.val, val) {
        return
    }
    {{- end }}
    *b.val = val
    b.trigger()
}

type boundExternal{{ .Name }} struct {
    bound{{ .Name }}
    old {{ .Type }}
}

func (b *boundExternal{{ .Name }}) Set(val {{ .Type }}) {
    b.lock.Lock()
    defer b.lock.Unlock()
    {{- if eq .Comparator "" }}
    if b.old == val {
        return
    }
    {{- else }}
    if {{ .Comparator }}(b.old, val) {
        return
    }
    {{- end }}
    *b.val = val
    b.old = val
    b.trigger()
}

func (b *boundExternal{{ .Name }}) Reload() {
    b.Set(*b.val)
}
`

type Sizeable interface {
	SizePtr() *geom.Point
}

type SizeEmbed struct {
	Size geom.Point
}

func (e *SizeEmbed) SizePtr() *geom.Point {
	return &e.Size
}



type embedValues struct {
	Name, Type, Interface string
}

func newFile(name string) (*os.File, error) {
	_, dirname, _, _ := runtime.Caller(0)
	filepath := path.Join(path.Dir(dirname), name+".go")
	os.Remove(filepath)
	f, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("Unable to open file %s: %v", f.Name(), err)
		return nil, err
	}

	f.WriteString(`// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package ledgrid
`)
	return f, nil
}

func writeFile(f *os.File, t *template.Template, d interface{}) {
	if err := t.Execute(f, d); err != nil {
		log.Fatalf("Unable to write file %s: %v", f.Name(), err)
	}
}

func main() {
	embedFile, err := newFile("animEmbed")
	if err != nil {
		return
	}
	defer embedFile.Close()
	embedFile.WriteString(`
import (
    "bytes"
)
`)

	embed := template.Must(template.New("embed").Parse(embedTemplate))
	binds := []embedValues{
		embedValues{Name: "Size", Type: "geom.Point", Interface: "Sizeable"},
		embedValues{Name: "Bytes", Type: "[]byte", Default: "nil", Comparator: "bytes.Equal"},
		embedValues{Name: "Float", Type: "float64", Default: "0.0", Format: "%f"},
		embedValues{Name: "Int", Type: "int", Default: "0", Format: "%d"},
		embedValues{Name: "Rune", Type: "rune", Default: "rune(0)"},
		embedValues{Name: "String", Type: "string", Default: "\"\""},
		//bindValues{Name: "Untyped", Type: "interface{}", Default: "nil"},
		//bindValues{Name: "Untyped", Type: "interface{}", Default: "nil", Since: "2.1"},
		//bindValues{Name: "URI", Type: "fyne.URI", Default: "fyne.URI(nil)", Since: "2.1",
		//FromString: "uriFromString", ToString: "uriToString", Comparator: "compareURI"},
	}
	for _, b := range binds {
		writeFile(itemFile, item, b)

		if b.Format != "" || b.ToString != "" {
			writeFile(convertFile, toString, b)
		}
	}
	for _, b := range binds {
		if b.Format != "" || b.ToString != "" {
			writeFile(convertFile, fromString, b)
		}
	}
}
