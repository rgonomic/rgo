// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codegen

import (
	"fmt"
	"go/types"
	"io"
	"path"
	"strings"
	"text/template"

	"github.com/rgonomic/rgo/internal/camel"
	"github.com/rgonomic/rgo/internal/pkg"
)

type FileSystem interface {
	Open(path string) (io.WriteCloser, error)
	Flush() error
}

func Render(dst FileSystem, pth string, tmpl *template.Template, info *pkg.Info) error {
	if len(info.Funcs) == 0 {
		return nil
	}

	if strings.Contains(pth, "%s") {
		pth = fmt.Sprintf(pth, path.Base(info.Pkg().Path()))
	}
	w, err := dst.Open(pth)
	if err != nil {
		return err
	}
	err = tmpl.Execute(w, info)
	if err != nil {
		return err
	}
	return w.Close()
}

// snake returns a closure that converts the input string camel case string
// to snake case which breaks also made based on the provided known words.
func snake(words []string) func(string) string {
	c := camel.NewSplitter(words)
	return func(s string) string {
		return strings.ToLower(strings.Join(c.Split(s), "_"))
	}
}

// names returns a comma-separated list of the names of the variables in vars.
func names(leadingComma bool, vars []*types.Var) string {
	if len(vars) == 0 {
		return ""
	}
	var buf strings.Builder
	for i, v := range vars {
		if leadingComma || i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(v.Name())
	}
	return buf.String()
}

// varsOf returns the vars of the given tuple.
func varsOf(t *types.Tuple) []*types.Var {
	if t.Len() == 0 {
		return nil
	}
	vars := make([]*types.Var, t.Len())
	for i := 0; i < t.Len(); i++ {
		vars[i] = t.At(i)
	}
	return vars
}
