// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codegen

import (
	"fmt"
	"go/types"
	"path"
	"reflect"
	"strings"
	"text/template"

	"github.com/rgonomic/rgo/internal/pkg"
)

// goFunc is the template for Go function file generation.
func GoFuncTemplate() *template.Template {
	return template.Must(template.New("Go func").Funcs(template.FuncMap{
		"imports":    imports,
		"varsOf":     varsOf,
		"go":         goParams,
		"anon":       anonymous,
		"types":      typeNames,
		"mangle":     pkg.Mangle,
		"unpackSEXP": unpackSEXPFuncGo,
		"packSEXP":   packSEXPFuncGo,
		"dec":        func(i int) int { return i - 1 },
	}).Parse(`{{$pkg := .Pkg}}// Code generated by rgnonomic/rgo; DO NOT EDIT.

package main

/*
#define USE_RINTERNALS
#include <R.h>
#include <Rinternals.h>
extern void R_error(char *s);

// TODO(kortschak): Only emit these when needed.
extern Rboolean Rf_isNull(SEXP s);
extern _GoString_ R_gostring(SEXP x, R_xlen_t i);
extern int getListElementIndex(SEXP list, const char *str);
*/
import "C"

import (
	"fmt"
	"unsafe"

{{with imports .}}{{range $p := .}}	"{{.}}"
{{end}}
{{end}}	"{{$pkg.Path}}"
)
{{$resultNeedsList := false}}{{range $func := .Funcs}}{{$params := varsOf $func.Signature.Params}}{{$results := varsOf $func.Signature.Results}}
//export Wrapped_{{$func.Name}}
func Wrapped_{{$func.Name}}({{go "_R_" $params}}) C.SEXP {
	defer func() {
		r := recover()
		if r != nil {
			err := C.CString(fmt.Sprint(r))
			C.R_error(err)
			C.free(unsafe.Pointer(err))
		}
	}()

	{{range $i, $p := $params}}_p{{$i}} := unpackSEXP{{mangle $p.Type}}(_R_{{$p.Name}})
	{{end}}{{with $results}}{{anon . "_r" false}} := {{end}}{{$pkg.Name}}.{{$func.Name}}({{anon $params "_p" false}}{{if $func.Signature.Variadic}}...{{end}})
	{{with $results}}return packSEXP_{{$func.Name}}({{anon . "_r" false}}){{else}}return C.R_NilValue{{end}}
}

{{if $results}}func packSEXP_{{$func.Name}}({{anon $results "p" true}}) C.SEXP {
{{$l := len $results -}}
{{- if eq $l 1 -}}
{{- $p := index $results 0}}	return packSEXP{{mangle $p.Type}}({{if $p.Name}}{{$p.Name}}{{else}}p0{{end -}})
{{- else}}{{$resultNeedsList = true}}	r := C.allocList({{len $results}})
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, {{len $results}})
	C.Rf_protect(names)
	arg := r
{{range $i, $p := $results}}{{$res := printf "r%d" $i}}{{if $p.Name}}{{$res = $p.Name}}{{end}}	C.SET_STRING_ELT(names, {{$i}}, C.Rf_mkCharLenCE(C._GoStringPtr("{{$res}}"), {{len $res}}, C.CE_UTF8))
	C.SETCAR(arg, packSEXP{{mangle $p.Type}}({{if $p.Name}}{{$p.Name}}{{else}}p{{$i}}{{end}}))
{{if lt $i (dec $l)}}	arg = C.CDR(arg)
{{end -}}
{{- end}}	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r{{end}}
}
{{end}}{{end}}
{{/* TODO(kortschak): Hoist C.SEXP unpacking for basic types out to the C code. */ -}}
{{- .Unpackers.Types | unpackSEXP -}}
{{- .Packers.Types | packSEXP}}func main() {}
`))
}

// goParams returns a comma-separated list of C.SEXP parameters using the
// parameter names in vars with the mangling prefix applied.
func goParams(prefix string, vars []*types.Var) string {
	if len(vars) == 0 {
		return ""
	}
	var buf strings.Builder
	for i, v := range vars {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(prefix)
		buf.WriteString(v.Name())
	}
	buf.WriteString(" C.SEXP")
	return buf.String()
}

// anonymous returns a comma-separated list of numbered parameters corresponding
// to vars with the given prefix. If typed is true, the parameter type is included.
func anonymous(vars []*types.Var, prefix string, typed bool) string {
	if len(vars) == 0 {
		return ""
	}
	var buf strings.Builder
	for i, v := range vars {
		if i != 0 {
			buf.WriteString(", ")
		}
		if !typed {
			buf.WriteString(fmt.Sprintf("%s%d", prefix, i))
			continue
		}
		name := v.Name()
		if name == "" {
			name = fmt.Sprintf("%s%d", prefix, i)
		}
		buf.WriteString(fmt.Sprintf("%s %s", name, path.Base(v.Type().String())))
	}
	return buf.String()
}

// typeNames returns a comma-separated list of the type names corresponding to vars.
func typeNames(vars []*types.Var) string {
	if len(vars) == 0 {
		return ""
	}
	var buf strings.Builder
	for i, v := range vars {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(nameOf(v.Type()))
	}
	return buf.String()
}

// nameOf returns the package name-qualified name of t.
func nameOf(t types.Type) string {
	return types.TypeString(t, func(pkg *types.Package) string {
		return pkg.Name()
	})
}

// targetFieldName returns the rgo struct tag of the ith field of s if
// it exists, otherwise the name of the field.
func targetFieldName(s *types.Struct, i int) string {
	tag := reflect.StructTag(s.Tag(i)).Get("rgo")
	if tag != "" {
		return tag
	}
	return s.Field(i).Name()
}
