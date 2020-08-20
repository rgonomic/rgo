// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestAnalyse(t *testing.T) {
	paths, err := filepath.Glob("./testdata/*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, path := range paths {
		fi, err := os.Stat(path)
		if err != nil {
			t.Errorf("failed to stat %q: %v", path, err)
			continue
		}
		if !fi.IsDir() {
			continue
		}

		info, err := Analyse(filepath.Join("github.com/rgonomic/rgo/internal/pkg", path), "", false)
		if err != nil {
			t.Errorf("unexpected error during analysis of %q: %v", path, err)
			continue
		}
		if info.Pkg().Name() != filepath.Base(path) {
			t.Errorf("unexpected package name: got:%q want:%q", info.Pkg().Name(), filepath.Base(path))
		}
		for _, fn := range info.Funcs {
			if !fn.Exported() {
				t.Errorf("unexpected unexported function: %s", fn)
			}
		}

		got := make(map[string][]string)
		for k, v := range info.Unpackers {
			got["in"] = append(got["in"], k)
			if k != v.String() {
				t.Errorf("unexpected type for unpacker key: %s != %s", k, v)
			}
		}
		for k, v := range info.Packers {
			got["out"] = append(got["out"], k)
			if k != v.String() {
				t.Errorf("unexpected type for packer key: %s != %s", k, v)
			}
		}
		for k := range got {
			sort.Strings(got[k])
		}

		want, err := typesFor(path)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("unexpected types in %q:\ngot: %v\nwant:%v", path, got, want)
		}
	}
}

func typesFor(path string) (map[string][]string, error) {
	cfg := &packages.Config{Mode: packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, "pattern=./"+path)
	if err != nil {
		return nil, fmt.Errorf("unexpected error loading %q: %w", path, err)
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("unexpected number of packages for %q: got:%d want:1", path, len(pkgs))
	}

	want := make(map[string][]string)
	for _, f := range pkgs[0].Syntax {
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Doc == nil {
				continue
			}

			var fn map[string][]string
			err := json.Unmarshal([]byte(strings.TrimPrefix(fd.Doc.List[0].Text, "//")), &fn)
			if err != nil {
				return nil, fmt.Errorf("failed to parse parameter types: %w", err)
			}
			for k, v := range fn {
				want[k] = append(want[k], v...)
			}
		}
	}
	for k := range want {
		t := want[k]
		sort.Strings(t)
		i := 0
		for _, v := range t {
			if v != t[i] {
				i++
				t[i] = v
			}
		}
		want[k] = t[:i+1]
	}

	return want, nil
}
