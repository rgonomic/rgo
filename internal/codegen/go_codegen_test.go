// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codegen

import (
	"bytes"
	"flag"
	"fmt"
	"go/types"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/diff"

	"github.com/rgonomic/rgo/internal/pkg"
)

var regenerate = flag.Bool("regen", false, "regenerate golden data from current state")

var mockPkg = types.NewPackage("path/to/pkg", "pkg")

var sexpFuncGoTests = []struct {
	typ types.Type
}{
	// Basic types.
	{typ: types.Typ[types.String]},
	{typ: types.Typ[types.Int32]},
	{typ: types.Universe.Lookup("rune").Type()},
	{typ: types.Typ[types.Uint8]},
	{typ: types.Universe.Lookup("byte").Type()},
	{typ: types.Typ[types.Float64]},
	{typ: types.Typ[types.Complex128]},
	{typ: types.Typ[types.Bool]},

	// Pointer types.
	{typ: types.NewPointer(types.Typ[types.String])},
	{typ: types.NewPointer(types.Typ[types.Int32])},
	{typ: types.NewPointer(types.Universe.Lookup("rune").Type())},
	{typ: types.NewPointer(types.Typ[types.Uint8])},
	{typ: types.NewPointer(types.Universe.Lookup("byte").Type())},
	{typ: types.NewPointer(types.Typ[types.Float64])},
	{typ: types.NewPointer(types.Typ[types.Complex128])},
	{typ: types.NewPointer(types.Typ[types.Bool])},

	// Array types.
	{typ: types.NewArray(types.Typ[types.String], 10)},
	{typ: types.NewArray(types.Typ[types.Int32], 10)},
	{typ: types.NewArray(types.Universe.Lookup("rune").Type(), 10)},
	{typ: types.NewArray(types.Typ[types.Uint8], 10)},
	{typ: types.NewArray(types.Universe.Lookup("byte").Type(), 10)},
	{typ: types.NewArray(types.Typ[types.Float64], 10)},
	{typ: types.NewArray(types.Typ[types.Complex128], 10)},
	{typ: types.NewArray(types.Typ[types.Bool], 10)},

	// Slice types.
	{typ: types.NewSlice(types.Typ[types.String])},
	{typ: types.NewSlice(types.Typ[types.Int32])},
	{typ: types.NewSlice(types.Universe.Lookup("rune").Type())},
	{typ: types.NewSlice(types.Typ[types.Uint8])},
	{typ: types.NewSlice(types.Universe.Lookup("byte").Type())},
	{typ: types.NewSlice(types.Typ[types.Float64])},
	{typ: types.NewSlice(types.Typ[types.Complex128])},
	{typ: types.NewSlice(types.Typ[types.Bool])},

	// Map types.
	{typ: types.NewMap(types.Typ[types.String], types.Typ[types.String])},
	{typ: types.NewMap(types.Typ[types.String], types.Typ[types.Int32])},
	{typ: types.NewMap(types.Typ[types.String], types.Universe.Lookup("rune").Type())},
	{typ: types.NewMap(types.Typ[types.String], types.Typ[types.Uint8])},
	{typ: types.NewMap(types.Typ[types.String], types.Universe.Lookup("byte").Type())},
	{typ: types.NewMap(types.Typ[types.String], types.Typ[types.Float64])},
	{typ: types.NewMap(types.Typ[types.String], types.Typ[types.Complex128])},
	{typ: types.NewMap(types.Typ[types.String], types.Typ[types.Bool])},

	// Struct types.
	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Typ[types.String], false),
			types.NewField(0, mockPkg, "F2", types.Typ[types.String], false),
		}, []string{`rgo:"Rname"`}),
	},

	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Typ[types.Int32], false),
			types.NewField(0, mockPkg, "F2", types.Typ[types.Int32], false),
		}, []string{`rgo:"Rname"`}),
	},
	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Universe.Lookup("rune").Type(), false),
			types.NewField(0, mockPkg, "F2", types.Universe.Lookup("rune").Type(), false),
		}, []string{`rgo:"Rname"`}),
	},

	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Typ[types.Uint8], false),
			types.NewField(0, mockPkg, "F2", types.Typ[types.Uint8], false),
		}, []string{`rgo:"Rname"`}),
	},
	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Universe.Lookup("byte").Type(), false),
			types.NewField(0, mockPkg, "F2", types.Universe.Lookup("byte").Type(), false),
		}, []string{`rgo:"Rname"`}),
	},

	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Typ[types.Float64], false),
			types.NewField(0, mockPkg, "F2", types.Typ[types.Float64], false),
		}, []string{`rgo:"Rname"`}),
	},

	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Typ[types.Complex128], false),
			types.NewField(0, mockPkg, "F2", types.Typ[types.Complex128], false),
		}, []string{`rgo:"Rname"`}),
	},

	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Typ[types.Bool], false),
			types.NewField(0, mockPkg, "F2", types.Typ[types.Bool], false),
		}, []string{`rgo:"Rname"`}),
	},
}

func TestUnpackSEXPFuncGo(t *testing.T) {
	if got := strings.TrimSpace(packSEXPFuncGo(nil)); got != "" {
		t.Errorf("unexpected output for empty slice: %s", got)
	}
	for i, test := range sexpFuncGoTests {
		typs := []types.Type{
			test.typ,
			types.NewNamed(types.NewTypeName(0, mockPkg, "T", nil), test.typ, nil),
		}
		for _, typ := range typs {
			got := []byte(strings.TrimSpace(unpackSEXPFuncGo([]types.Type{typ})))

			var named string
			if _, ok := typ.(*types.Named); ok {
				named = "-named"
			}
			golden := filepath.Join("testdata", fmt.Sprintf("unpackSEXP%s%s.golden", pkg.Mangle(test.typ), named))
			if *regenerate {
				err := ioutil.WriteFile(golden, got, 0o664)
				if err != nil {
					t.Fatalf("failed to write golden data: %v", err)
				}
				continue
			}

			want, err := ioutil.ReadFile(golden)
			if err != nil {
				t.Fatalf("failed to read golden data: %v", err)
			}

			if !bytes.Equal(got, want) {
				var buf bytes.Buffer
				err := diff.Text("got", "want", got, want, &buf)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				t.Errorf("unexpected generated code for test %d %s (unpackSEXP%s.golden):\n%s", i, typ, pkg.Mangle(typ), &buf)
			}
		}
	}
}

func TestPackSEXPFuncGo(t *testing.T) {
	if got := strings.TrimSpace(packSEXPFuncGo(nil)); got != "" {
		t.Errorf("unexpected output for empty slice: %s", got)
	}
	for i, test := range sexpFuncGoTests {
		typs := []types.Type{
			test.typ,
			types.NewNamed(types.NewTypeName(0, mockPkg, "T", nil), test.typ, nil),
		}
		for _, typ := range typs {
			got := []byte(strings.TrimSpace(packSEXPFuncGo([]types.Type{typ})))

			var named string
			if _, ok := typ.(*types.Named); ok {
				named = "-named"
			}
			golden := filepath.Join("testdata", fmt.Sprintf("packSEXP%s%s.golden", pkg.Mangle(test.typ), named))
			if *regenerate {
				err := ioutil.WriteFile(golden, got, 0o664)
				if err != nil {
					t.Fatalf("failed to write golden data: %v", err)
				}
				continue
			}

			want, err := ioutil.ReadFile(golden)
			if err != nil {
				t.Fatalf("failed to read golden data: %v", err)
			}

			if !bytes.Equal(got, want) {
				var buf bytes.Buffer
				err := diff.Text("got", "want", got, want, &buf)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				t.Errorf("unexpected generated code for test %d %s:\n%s", i, golden, &buf)
			}
		}
	}
}
