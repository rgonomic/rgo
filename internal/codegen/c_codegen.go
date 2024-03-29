// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codegen

import (
	"go/types"
	"strings"
	"text/template"
)

// cFunc is the template for C shim function file generation.
func CFuncTemplate(words []string) *template.Template {
	return template.Must(template.New("C func").Funcs(template.FuncMap{
		"snake":  snake(words),
		"varsOf": varsOf,
		"c":      cParams,
		"names":  names,
	}).Parse(`// Code generated by rgnonomic/rgo; DO NOT EDIT.

#include "_cgo_export.h"

void R_warning(char* s) {
	warning(s);
}

void R_error(char* s) {
	error(s);
}

// TODO(kortschak): Only emit these when needed:
// Needed for unpacking SEXP character.
GoString R_gostring(SEXP x, R_xlen_t i) {
	SEXP _s = STRING_ELT(x, i);
	GoString s = {(char*)CHAR(_s), LENGTH(_s)};
	return s;
}

// Needed for getting list elements by name.
int getListElementIndex(SEXP list, const char *str) {
	int index = -1;
	SEXP names = getAttrib(list, R_NamesSymbol);
	if (!isString(names)) {
		return index;
	}
	for (int i = 0; i < length(list); i++) {
		if (strcmp(CHAR(STRING_ELT(names, i)), str) == 0) {
			index = i;
			break;
		}
	}
	return index;
}{{range $func := .Funcs}}{{$params := varsOf $func.Signature.Params}}

SEXP {{snake $func.Func.Name}}({{c $params}}) {
	return Wrapped_{{$func.Func.Name}}({{names false $params}});
}{{end}}
`))
}

// cParams returns a comma-separated list of the variable names in vars
// a set of C SEXP parameters.
func cParams(vars []*types.Var) string {
	if len(vars) == 0 {
		return ""
	}
	var buf strings.Builder
	for i, v := range vars {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("SEXP ")
		buf.WriteString(v.Name())
	}
	return buf.String()
}
