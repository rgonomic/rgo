// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build generate

package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	for _, g := range generators {
		f, err := os.Create(g.file)
		if err != nil {
			log.Fatal(err)
		}
		err = g.template.Execute(f, g.details)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
	}
}

var generators = []struct {
	file     string
	template *template.Template
	details  interface{}
}{
	{file: "vector.go", template: vectorTemplate, details: vectors},
	{file: "list.go", template: listTemplate, details: lists},
	{file: "funcs.go", template: funcsTemplate, details: funcs},
}

var vectors = []struct {
	Type      string
	TypeDoc   string
	NewDoc    string
	Allocator string
	RType     string
	RTypeSpec string
	ElemType  string
	VectorDoc string
	Stringer  string
}{
	{
		Type: "Integer", TypeDoc: "an R integer vector",
		NewDoc: "an integer vector with length n", Allocator: "Vector",
		RType: "vector_sexprec", RTypeSpec: "INTSXP",
		ElemType: "int32",
	},
	{
		Type: "Logical", TypeDoc: "an R logical vector",
		NewDoc: "a logical vector with length n", Allocator: "Vector",
		RType: "vector_sexprec", RTypeSpec: "LGLSXP",
		ElemType: "int32",
	},
	{
		Type: "Real", TypeDoc: "an R real vector",
		NewDoc: "a real vector with length n", Allocator: "Vector",
		RType: "vector_sexprec", RTypeSpec: "REALSXP",
		ElemType: "float64",
	},
	{
		Type: "Complex", TypeDoc: "an R complex vector",
		NewDoc: "a complex vector with length n", Allocator: "Vector",
		RType: "vector_sexprec", RTypeSpec: "CPLXSXP",
		ElemType: "complex128",
	},
	{
		Type: "String", TypeDoc: "an R character vector",
		NewDoc: "a character vector with length n", Allocator: "Vector",
		RType: "vector_sexprec", RTypeSpec: "STRSXP",
		ElemType: "*Character",
	},
	{
		Type: "Character", TypeDoc: "the R representation of a string",
		NewDoc: "a scalar string corresponding to s", Allocator: "String",
		RType: "vector_sexprec", RTypeSpec: "STRSXP",
		ElemType: "byte",
		Stringer: `// String returns a Go string corresponding the the R characters.
// The returned string is allocated by the Go runtime.
func (v *Character) String() string {
	return string(v.Bytes())
}`,
		VectorDoc: "Bytes returns the bytes held by the R SEXP value",
	},
	{
		Type: "Raw", TypeDoc: "an R raw vector",
		NewDoc: "a raw vector with length n", Allocator: "Vector",
		RType: "vector_sexprec", RTypeSpec: "RAWSXP",
		ElemType:  "byte",
		VectorDoc: "Bytes returns the bytes held by the R SEXP value",
	},
	{
		Type: "Vector", TypeDoc: "a generic R vector",
		NewDoc: "a generic vector with length n", Allocator: "Vector",
		RType: "vector_sexprec", RTypeSpec: "VECSXP",
		ElemType: "*Value",
	},
	{
		Type: "Expression", TypeDoc: "an R expression",
		RType:     "vector_sexprec",
		ElemType:  "*Value",
		VectorDoc: "Vector returns a slice corresponding to the R expression",
	},
	{
		Type: "WeakReference", TypeDoc: "an R weak reference",
		RType:     "vector_sexprec",
		ElemType:  "*Value",
		VectorDoc: "Vector returns the four elements of an R weak reference",
	},
}

var vectorTemplate = template.Must(template.New("vectors").Parse(`// Code generated by 'go generate github.com/rgonomic/rg/sexp'; DO NOT EDIT.

// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sexp

import "unsafe"

// base returns the address of the first element of the vector.
func (v *vector_sexprec) base() unsafe.Pointer {
	return add(unsafe.Pointer(v), unsafe.Sizeof(vector_sexprec{}))
}

func add(addr unsafe.Pointer, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(addr) + offset)
}
{{range .}}
// {{.Type}} is {{.TypeDoc}}.
type {{.Type}} struct {
	{{.RType}}
}

{{if .Allocator}}// New{{.Type}} returns {{.NewDoc}}.
//
// The allocation is made by the R runtime. The returned value may need to
// call its Protect method.
func New{{.Type}}({{if eq .Allocator "String"}}s string{{else}}n int{{end}}) *{{.Type}} {
	return (*{{.Type}})(allocate{{.Allocator}}({{if eq .Allocator "String"}}s{{else}}{{.RTypeSpec}}, n{{end}}))
}

// Protect protects the SEXP value and returns it.
func (v *{{.Type}}) Protect() *{{.Type}} {
	if v == nil {
		return nil
	}
	return (*{{.Type}})(protect(unsafe.Pointer(v)))
}

// Unprotect unprotects the SEXP. It is equivalent to UNPROTECT(1).
func (v *{{.Type}}) Unprotect() {
	if v == nil || v.Value().IsNull() {
		panic("sexp: unprotecting a nil value")
	}
	unprotect(1)
}

{{end}}// Len returns the number of elements in the vector.
func (v *{{.Type}}) Len() int {
	if v == nil {
		return 0
	}
	return int(v.vecsxp.length)
}

// Info returns the information field of the SEXP value.
func (v *{{.Type}}) Info() Info {
	if v == nil {
		return NilValue.Info()
	}
	return *(*Info)(unsafe.Pointer(&v.sxpinfo))
}

// Value returns the generic state of the SEXP value.
func (v *{{.Type}}) Value() *Value {
	return (*Value)(unsafe.Pointer(v))
}

// Attributes returns the attributes of the SEXP value.
func (v *{{.Type}}) Attributes() *List {
	if v == nil {
		return nil
	}
	attr := (*List)(unsafe.Pointer(v.attrib))
	if attr.Value().IsNull() {
		return nil
	}
	return attr
}

// {{with .VectorDoc}}{{.}}{{else}}Vector returns a slice corresponding to the R vector{{end}}.
func (v *{{.Type}}) {{if eq .ElemType "byte"}}Bytes{{else}}Vector{{end}}() []{{.ElemType}} {
	if v == nil {
		return nil
	}
	n := v.Len()
	return (*[1 << 30]{{.ElemType}})(v.base())[:n:n]
}
{{with .Stringer}}
{{.}}
{{end}}{{end}}`))

var lists = []struct {
	Type      string
	TypeDoc   string
	NewDoc    string
	Allocator string
	RType     string
	HasGetter bool
	TagDoc    string
}{
	{
		Type: "List", TypeDoc: "an R linked list",
		NewDoc: "a list with length n", Allocator: "List",
		RType:     "list_sexprec",
		HasGetter: true,
		TagDoc:    "the list's tag value",
	},
	{
		Type: "Lang", TypeDoc: "an R language object",
		RType:  "list_sexprec",
		TagDoc: "the object's tag value",
	},
	{
		Type: "Dot", TypeDoc: "an R pairlist of promises",
		RType:  "list_sexprec",
		TagDoc: "the object's tag value",
	},
	{
		Type: "Symbol", TypeDoc: "an R name value",
		RType: "sym_sexprec",
	},
}

var listTemplate = template.Must(template.New("lists").Parse(`// Code generated by 'go generate github.com/rgonomic/rg/sexp'; DO NOT EDIT.

// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sexp

import (
	"bytes"
	"unsafe"
)
{{range .}}
// {{.Type}} is {{.TypeDoc}}.
type {{.Type}} struct {
	{{.RType}}
}

{{if .Allocator}}// New{{.Type}} returns {{.NewDoc}}.
func New{{.Type}}(n int) *{{.Type}} {
	return (*{{.Type}})(allocate{{.Allocator}}(n))
}

// Protect protects the SEXP value and returns it.
func (v *{{.Type}}) Protect() *{{.Type}} {
	if v == nil {
		return nil
	}
	return (*{{.Type}})(protect(unsafe.Pointer(v)))
}

// Unprotect unprotects the SEXP. It is equivalent to UNPROTECT(1).
func (v *{{.Type}}) Unprotect() {
	if v == nil || v.Value().IsNull() {
		panic("sexp: unprotecting a nil value")
	}
	unprotect(1)
}

{{end}}// Info returns the information field of the SEXP value.
func (v *{{.Type}}) Info() Info {
	if v == nil {
		return NilValue.Info()
	}
	return *(*Info)(unsafe.Pointer(&v.sxpinfo))
}

// Value returns the generic state of the SEXP value.
func (v *{{.Type}}) Value() *Value {
	return (*Value)(unsafe.Pointer(v))
}

// Attributes returns the attributes of the SEXP value.
func (v *{{.Type}}) Attributes() *List {
	if v == nil {
		return nil
	}
	attr := (*List)(unsafe.Pointer(v.attrib))
	if attr.Value().IsNull() {
		return nil
	}
	return attr
}

{{if eq .RType "list_sexprec"}}// Head returns the first element of the list (CAR/lisp in R terminology).
func (v *{{.Type}}) Head() *Value {
	if v == nil || v.Value().IsNull() {
		return nil
	}
	car := (*Value)(unsafe.Pointer(v.list_sxp.carval))
	if car.IsNull() {
		return nil
	}
	return car
}

// Tail returns the remaining elements of the list (CDR/lisp in R terminology).
func (v *{{.Type}}) Tail() *Value {
	if v == nil || v.Value().IsNull() {
		return nil
	}
	cdr := (*Value)(unsafe.Pointer(v.list_sxp.cdrval))
	if cdr.IsNull() {
		return nil
	}
	return cdr
}

// Tag returns {{.TagDoc}}.
func (v *{{.Type}}) Tag() *Symbol {
	if v == nil || v.Value().IsNull() {
		return nil
	}
	tag := (*Symbol)(unsafe.Pointer(v.list_sxp.tagval))
	if tag.Value().IsNull() {
		return nil
	}
	return tag
}
{{if .HasGetter}}
// Get returns the the Value associated with the given tag in the list.
func (v *{{.Type}}) Get(tag []byte) *Value {
	curr := v
	for !curr.Value().IsNull() {
		t := curr.Tag()
		if t != nil {
			if bytes.Equal(t.Name().Bytes(), tag) {
				return (*Value)(curr.Pointer())
			}
		}
		tail := curr.Tail()
		if tail, ok := tail.Value().Interface().(*{{.Type}}); ok {
			curr = (*{{.Type}})(tail.Pointer())
			continue
		}
		break
	}
	return nil
}

// tags returns all the tags for the list. The []string is allocated
// by the Go runtime.
func (v *{{.Type}}) tags() []string {
	var tags []string
	curr := v
	for !curr.Value().IsNull() {
		t := curr.Tag()
		if t != nil {
			tag := t.String()
			tags = append(tags, tag)
		}
		tail := curr.Tail()
		if tail, ok := tail.Value().Interface().(*{{.Type}}); ok {
			curr = (*{{.Type}})(tail.Pointer())
			continue
		}
		break
	}
	return tags
}
{{end}}{{end}}{{if eq .RType "sym_sexprec"}}// Value returns the value of the symbol.
func (v *{{.Type}}) SymbolValue() *Value {
	if v == nil {
		return nil
	}
	val := (*Value)(unsafe.Pointer(v.sym_sxp.value))
	if val.IsNull() {
		return nil
	}
	return val
}

// Name returns the name of the symbol
func (v *{{.Type}}) Name() *Character {
	if v == nil {
		return nil
	}
	name := (*Character)(unsafe.Pointer(v.sym_sxp.pname))
	if name.Value().IsNull() {
		return nil
	}
	return name
}

// String returns a Go string of the symbol name.
// The returned string is allocated by the Go runtime.
func (v *{{.Type}}) String() string {
	return v.Name().String()
}

// Internal returns a pointer if the symbol is a .Internal function.
func (v *{{.Type}}) Internal() *Value {
	if v == nil {
		return nil
	}
	intern := (*Value)(unsafe.Pointer(v.sym_sxp.internal))
	if intern.IsNull() {
		return nil
	}
	return intern
}{{end}}{{end}}
`))

type Accessor struct {
	Name  string
	Doc   string
	Field string
}

var funcs = []struct {
	Type      string
	TypeDoc   string
	RType     string
	Member    string
	Accessors []Accessor
}{
	{
		Type: "Promise", TypeDoc: "an R promise",
		RType: "prom_sexprec", Member: "prom_sxp",
		Accessors: []Accessor{
			{Name: "PromiseValue", Doc: "the value of the promise", Field: "value"},
			{Name: "Expression", Doc: "the expression to be evaluated", Field: "expr"},
			{Name: "Environment", Doc: "the environment in which to evaluate the expression", Field: "env"},
		},
	},
	{
		Type: "Closure", TypeDoc: "an R closure",
		RType: "clo_sexprec", Member: "clos_sxp",
		Accessors: []Accessor{
			{Name: "Formals", Doc: "the formal arguments of the function", Field: "formals"},
			{Name: "Body", Doc: "the body of the function", Field: "body"},
			{Name: "Environment", Doc: "the environment in which to evaluate the function", Field: "env"},
		},
	},
	{
		Type: "Environment", TypeDoc: "a current execution environment",
		RType: "env_sexprec", Member: "env_sxp",
		Accessors: []Accessor{
			{Name: "Frame", Doc: "the current frame", Field: "frame"},
			{Name: "Enclosing", Doc: "the enclosing environment", Field: "enclos"},
			{Name: "HashTable", Doc: "the environment's hash table", Field: "hashtab"},
		},
	},
	{
		Type: "Builtin", TypeDoc: "an R language built-in function",
		RType: "prim_sexprec",
	},
	{
		Type: "Special", TypeDoc: "an R language built-in function",
		RType: "prim_sexprec",
	},
}

var funcsTemplate = template.Must(template.New("funcs").Parse(`// Code generated by 'go generate github.com/rgonomic/rg/sexp'; DO NOT EDIT.

// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sexp

import "unsafe"
{{range .}}{{$t := .}}
// {{.Type}} is {{.TypeDoc}}.
type {{.Type}} struct {
	{{.RType}}
}

// Info returns the information field of the SEXP value.
func (v *{{.Type}}) Info() Info {
	if v == nil {
		return NilValue.Info()
	}
	return *(*Info)(unsafe.Pointer(&v.sxpinfo))
}

// Value returns the generic state of the SEXP value.
func (v *{{.Type}}) Value() *Value {
	return (*Value)(unsafe.Pointer(v))
}

// Attributes returns the attributes of the SEXP value.
func (v *{{.Type}}) Attributes() *List {
	if v == nil {
		return nil
	}
	attr := (*List)(unsafe.Pointer(v.attrib))
	if attr.Value().IsNull() {
		return nil
	}
	return attr
}
{{with .Accessors}}{{range .}}
// {{.Name}} returns {{.Doc}}.
func (v *{{$t.Type}}) {{.Name}}() *Value {
	if v == nil || v.Value().IsNull() {
		return nil
	}
	{{.Field}} := (*Value)(unsafe.Pointer(v.{{$t.Member}}.{{.Field}}))
	if {{.Field}}.IsNull() {
		return nil
	}
	return {{.Field}}
}
{{end}}{{else}}
// Offset returns the offset into the table of language primitives.
func (v *{{.Type}}) Offset() int32 {
	return int32(v.prim_sxp.offset)
}
{{end}}{{end}}`))
