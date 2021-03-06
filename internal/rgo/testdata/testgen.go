// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run ./testgen.go

// The testgen command constructs small single file packages with no-op, but
// buildable source for testing the internal/pkg code.
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

var crafted = []pkg{
	{
		Name:  "mixed",
		Path:  "github.com/rgonomic/rgo/internal/rgo/testdata",
		Types: []string{"T int", "S1 string"},
		Funcs: []fn{
			{In: []string{"int"}, Out: []string{"int", "int"}, Named: false},
			{In: []string{"int"}, Out: []string{"float64", "int"}, Named: true},
			{In: []string{"int"}},
			{In: []string{"T", "S1"}, Out: []string{"S1"}},
		},
	},
	{
		Name: "slice_of_slices",
		Path: "github.com/rgonomic/rgo/internal/rgo/testdata",
		Funcs: []fn{
			{In: []string{"[][]float64"}, Out: []string{"[][]float64"}, Named: false},
		},
	},
	{
		Name: "map_of_slices",
		Path: "github.com/rgonomic/rgo/internal/rgo/testdata",
		Funcs: []fn{
			{In: []string{"map[string][]float64"}, Out: []string{"map[string][]float64"}, Named: false},
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
	pkg   *pkg
	In    []string // Input parameter types.
	Out   []string // Output parameter types.
	Named bool     // Whether the output parameter types are named.
}

func builtins() []pkg {
	skip := map[types.BasicKind]bool{
		types.Uintptr: true,
		types.Uint64:  true,
		types.Int64:   true,
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

	// Generate scalar value test functions.
	dst = append(dst,
		pkg{Name: fmt.Sprintf("%s_in", name), Funcs: []fn{
			{In: []string{name}}}},
		pkg{Name: fmt.Sprintf("%s_out", name), Funcs: []fn{
			{Out: []string{name}}}},
		pkg{Name: fmt.Sprintf("%s_out_named", name), Funcs: []fn{
			{Out: []string{name}, Named: true}}},
	)

	dst = append(dst,
		pkg{Name: fmt.Sprintf("%s_slice_in", name), Funcs: []fn{
			{In: []string{"[]" + name}}}},
		pkg{Name: fmt.Sprintf("%s_slice_out", name), Funcs: []fn{
			{Out: []string{"[]" + name}}}},
		pkg{Name: fmt.Sprintf("%s_slice_out_named", name), Funcs: []fn{
			{Out: []string{"[]" + name}, Named: true}}},
	)

	// Generate array value test functions.
	dst = append(dst,
		pkg{Name: fmt.Sprintf("%s_array_in", name), Funcs: []fn{
			{In: []string{"[4]" + name}}}},
		pkg{Name: fmt.Sprintf("%s_array_out", name), Funcs: []fn{
			{Out: []string{"[4]" + name}}}},
		pkg{Name: fmt.Sprintf("%s_array_out_named", name), Funcs: []fn{
			{Out: []string{"[4]" + name}, Named: true}}},
	)

	// Generate struct value test functions.
	st := fmt.Sprintf(`struct{F1 %[1]s; F2 %[1]s "rgo:\"Rname\""}`, name)
	dst = append(dst,
		pkg{Name: fmt.Sprintf("struct_%s_in", name), Funcs: []fn{
			{In: []string{st}}}},
		pkg{Name: fmt.Sprintf("struct_%s_out", name), Funcs: []fn{
			{Out: []string{st}}}},
		pkg{Name: fmt.Sprintf("struct_%s_out_named", name), Funcs: []fn{
			{Out: []string{st}, Named: true}}},
	)

	// Generate map[string]T value test functions.
	mt := fmt.Sprintf("map[string]%s", name)
	dst = append(dst,
		pkg{Name: fmt.Sprintf("string_%s_map_in", name), Funcs: []fn{
			{In: []string{mt}}}},
		pkg{Name: fmt.Sprintf("string_%s_map_out", name), Funcs: []fn{
			{Out: []string{mt}}}},
		pkg{Name: fmt.Sprintf("string_%s_map_out_named", name), Funcs: []fn{
			{Out: []string{mt}, Named: true}}},
	)

	return dst
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

			err = os.Remove(filepath.Join(pkg, "go.mod"))
			if err != nil && !os.IsNotExist(err) {
				log.Fatalf("failed to remove go.mod file: %v", err)
			}
			cmd := exec.Command("go", "mod", "init", pkg)
			cmd.Dir = filepath.Join(".", pkg)
			err = cmd.Run()
			if err != nil {
				log.Fatalf("failed create go.mod file: %v", err)
			}
		}
	}
}

var src = template.Must(template.New("Go source").Parse(`// Code generated by "go generate github.com/rgonomic/rgo/internal/pkg/testdata"; DO NOT EDIT.

package {{.Name}}_{{.UID}}
{{if .Types}}
type (
{{- range $i, $t := .Types}}
	{{$t}}{{end}}
){{end}}
{{- range $i, $fn := .Funcs}}

// Test{{$i}} does things with {{$fn.In}} and returns {{$fn.Out}}.
func Test{{$i}}({{if $fn.In}}{{range $j, $p := $fn.In -}}
	{{- if ne $j 0}}, {{end}}par{{$j}} {{$p -}}
{{- end}}{{end}}){{if $fn.Out}} ({{range $j, $p := $fn.Out -}}
	{{- if ne $j 0}}, {{end}}{{if $fn.Named}}res{{$j}} {{end}}{{$p}}{{end}}){{end}} { {{if not $fn.Named}}{{range $j, $p := $fn.Out}}
	var res{{$j}} {{$p}}{{end}}{{end}}
{{if $fn.Out}}	return {{range $j, $p := $fn.Out -}}
	{{- if ne $j 0}}, {{end}}res{{$j}}{{end}}
{{end}}}{{end}}
`))
