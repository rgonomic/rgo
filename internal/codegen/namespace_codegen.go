// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codegen

import (
	"text/template"
)

func NamespaceTemplate(words []string) *template.Template {
	return template.Must(template.New("NAMESPACE").Funcs(template.FuncMap{
		"snake": snake(words),
	}).Parse(`# Code generated by rgnonomic/rgo; DO NOT EDIT.

useDynLib({{$.Pkg.Name}})
{{range $func := .Funcs}}export({{snake $func.Func.Name}})
{{end}}`))
}
