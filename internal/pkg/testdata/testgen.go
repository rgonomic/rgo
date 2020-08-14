// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run ./testgen.go

// The testgen command constructs small single file packages with no-op, but
// buildable source for testing the internal/pkg code.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"go/token"
	"go/types"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"text/template"
)

var crafted = []pkg{
	{
		Name:  "mixed",
		Path:  "github.com/rgonomic/rgo/internal/pkg/testdata",
		Types: []string{"T int", "S1 string"},
		Funcs: []fn{
			{In: []string{"int"}, Out: []string{"int", "int"}, Named: false},
			{In: []string{"int"}, Out: []string{"float64", "int"}, Named: true},
			{In: []string{"int"}},
			{In: []string{"T", "S1"}, HelpIn: []string{"int", "string"}, Out: []string{"S1"}, HelpOut: []string{"string"}},
		},
	},
}

type pkg struct {
	Name  string
	UID   int
	Path  string
	Types []string
	Funcs []fn
}

type fn struct {
	pkg     *pkg
	In      []string `json:"in,omitempty"`  // Input parameter types.
	HelpIn  []string `json:"-"`             // Composing types for input parameters.
	Out     []string `json:"out,omitempty"` // Output parameter types.
	HelpOut []string `json:"-"`             // Composing types for output parameters.
	Named   bool     `json:"-"`             // Whether the output parameter types are named.
}

func builtins() []pkg {
	skip := map[types.BasicKind]bool{
		types.Uint64: true,
		types.Int64:  true,
	}
	var pkgs []pkg
	for t := types.Bool; t <= types.String; t++ {
		if skip[t] {
			continue
		}
		pkgs = addTypeTest(pkgs, types.Typ[t])
	}
	pkgs = addTypeTest(pkgs, types.Universe.Lookup("byte").Type().(*types.Basic))
	pkgs = addTypeTest(pkgs, types.Universe.Lookup("rune").Type().(*types.Basic))

	return pkgs
}

func addTypeTest(dst []pkg, typ *types.Basic) []pkg {
	name := typ.String()
	mapHelpIn := []string{"string", name}
	mapHelpOut := []string{"string", name}
	var cmplxHelp []string
	if typ.Kind() == types.Complex64 {
		cmplxHelp = []string{"complex128"}
		mapHelpIn = append(mapHelpIn, "complex128")
	}
	return append(dst,
		pkg{Name: fmt.Sprintf("%s_in", name), Funcs: []fn{
			{In: []string{name}, HelpIn: cmplxHelp}}},
		pkg{Name: fmt.Sprintf("%s_out", name), Funcs: []fn{
			{Out: []string{name}}}},
		pkg{Name: fmt.Sprintf("%s_out_named", name), Funcs: []fn{
			{Out: []string{name}, Named: true}}},

		pkg{Name: fmt.Sprintf("%s_slice_in", name), Funcs: []fn{
			{In: []string{"[]" + name}}}},
		pkg{Name: fmt.Sprintf("%s_slice_out", name), Funcs: []fn{
			{Out: []string{"[]" + name}}}},
		pkg{Name: fmt.Sprintf("%s_slice_out_named", name), Funcs: []fn{
			{Out: []string{"[]" + name}, Named: true}}},

		pkg{Name: fmt.Sprintf("string_%s_map_in", name), Funcs: []fn{
			{In: []string{fmt.Sprintf("map[string]%s", name)}, HelpIn: mapHelpIn}}},
		pkg{Name: fmt.Sprintf("string_%s_map_out", name), Funcs: []fn{
			{Out: []string{fmt.Sprintf("map[string]%s", name)}, HelpOut: mapHelpOut}}},
		pkg{Name: fmt.Sprintf("string_%s_map_out_named", name), Funcs: []fn{
			{Out: []string{fmt.Sprintf("map[string]%s", name)}, HelpOut: mapHelpOut, Named: true}}},
	)

}

func main() {
	suff := make(map[string]int)
	for _, cases := range [][]pkg{builtins(), crafted} {
		for _, c := range cases {
			c.UID = suff[c.Name]
			suff[c.Name]++
			for i := range c.Funcs {
				c.Funcs[i].pkg = &c
			}

			pkg := fmt.Sprintf("%s_%d", c.Name, c.UID)
			err := os.Mkdir(pkg, 0o755)
			if err != nil && !os.IsExist(err) {
				log.Fatalf("failed to create testing package dir: %v", err)
			}
			f, err := os.Create(filepath.Join(pkg, c.Name+".go"))
			if err != nil {
				log.Fatalf("failed to create testing source file: %v", err)
			}

			var buf bytes.Buffer
			err = src.Execute(&buf, c)
			if err != nil {
				log.Fatalf("failed to execute template: %v", err)
			}
			b, err := format.Source(buf.Bytes())
			if err != nil {
				log.Fatalf("failed to format source: %v", err)
			}
			_, err = f.Write(b)
			if err != nil {
				log.Fatalf("failed to write source: %v", err)
			}

			err = f.Close()
			if err != nil {
				log.Fatalf("failed to close testing source file: %v", err)
			}
		}
	}
}

func (f fn) JSON() string {
	pth := path.Join(f.pkg.Path, fmt.Sprintf("%s_%d", f.pkg.Name, f.pkg.UID))
	var t fn
	t.In = uniq(f.In, f.HelpIn...)
	addPathPrefix(pth, t.In)
	t.Out = uniq(f.Out, f.HelpOut...)
	addPathPrefix(pth, t.Out)
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func addPathPrefix(path string, types []string) {
	for i, typ := range types {
		if token.IsExported(typ) {
			types[i] = fmt.Sprintf("%s.%s", path, typ)
		}
	}
}

func uniq(s []string, extra ...string) []string {
	if len(s)+len(extra) == 0 {
		return nil
	}
	t := make([]string, len(s), len(s)+len(extra))
	copy(t, s)
	t = append(t, extra...)

	// Convert basic types to their kind.
	for i, v := range t {
		typ := types.Universe.Lookup(v)
		if typ == nil {
			continue
		}
		switch typ := typ.Type().(type) {
		case *types.Basic:
			t[i] = types.Typ[typ.Kind()].String()
		}
	}

	sort.Strings(t)
	i := 0
	for _, v := range t {
		if v != t[i] {
			i++
			t[i] = v
		}
	}
	return t[:i+1]
}

var src = template.Must(template.New("Go source").Parse(`// Code generated by "go generate github.com/rgonomic/rgo/internal/pkg/testdata"; DO NOT EDIT.

package {{.Name}}_{{.UID}}
{{if .Types}}
type (
{{- range $i, $t := .Types}}
	{{$t}}{{end}}
){{end}}
{{- range $i, $fn := .Funcs}}

//{{.JSON}}
func Test{{$i}}({{if $fn.In}}{{range $j, $p := $fn.In -}}
	{{- if ne $j 0}}, {{end}}par{{$j}} {{$p -}}
{{- end}}{{end}}){{if $fn.Out}} ({{range $j, $p := $fn.Out -}}
	{{- if ne $j 0}}, {{end}}{{if $fn.Named}}res{{$j}} {{end}}{{$p}}{{end}}){{end}} { {{if not $fn.Named}}{{range $j, $p := $fn.Out}}
	var res{{$j}} {{$p}}{{end}}{{end}}
{{if $fn.Out}}	return {{range $j, $p := $fn.Out -}}
	{{- if ne $j 0}}, {{end}}res{{$j}}{{end}}
{{end}}}{{end}}
`))
