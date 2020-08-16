// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"os/exec"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/tools/go/packages"
)

// Module returns module information for the module containing the given path.
func Module(path string) (ModInfo, error) {
	args := []string{"list", "-json", "-m"}
	if path != "" && path != "." {
		args = append(args, path)
	}
	var m ModInfo
	cmd := exec.Command("go", args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(buf.Bytes(), &m)
	return m, err
}

// ModInfo is information about a module. The values held by this struct
// are described in the output from `go help list`.
type ModInfo struct {
	Path    string // module path
	Version string // module version
	Dir     string // directory holding files for this module, if any
	GoMod   string // path to go.mod file used when loading this module, if any
}

// Info holds information about the functions in a package.
type Info struct {
	Funcs []FuncInfo

	Unpackers unpackers
	Packers   packers
}

func (p *Info) Pkg() *types.Package {
	if len(p.Funcs) == 0 {
		return nil
	}
	return p.Funcs[0].Pkg()
}

// FuncInfo holds type and syntax information about a function.
type FuncInfo struct {
	*types.Func
	*ast.FuncDecl
}

func (f FuncInfo) Signature() *types.Signature {
	return f.Func.Type().(*types.Signature)
}

func Analyse(path string, verbose bool) (*Info, error) {
	if strings.HasSuffix(path, "...") {
		return nil, errors.New("pkg: invalid use of ... suffix")
	}

	cfg := &packages.Config{
		Mode: packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
	}

	pkgs, err := packages.Load(cfg, "pattern="+path)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) != 0 {
		return nil, errors.New("package errors")
	}
	var pkg *packages.Package
	switch len(pkgs) {
	case 0:
		return nil, errors.New("pkg: no package analysed")
	case 1:
		pkg = pkgs[0]
	default:
		return nil, errors.New("pkg: more than one package analysed")
	}

	log.Printf("wrapping: %s", pkg.ID)
	if verbose {
		log.Println("files:", pkg.GoFiles)
	}
	var funcs []FuncInfo
	needUnpack := make(unpackers)
	needPack := make(packers)
	for _, f := range pkg.Syntax {
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			fn := pkg.TypesInfo.Defs[fd.Name].(*types.Func)
			if !fn.Exported() {
				if verbose {
					log.Printf("skipping %s: unexported function", fn.Name())
				}
				continue
			}
			sig := fn.Type().(*types.Signature)
			if sig.Recv() != nil {
				if verbose {
					log.Printf("skipping %s.%s: method function", sig.Recv().Type(), fn.Name())
				}
				continue
			}

			par := sig.Params()
			err := checkType(par, par, false)
			if err != nil {
				if verbose {
					log.Printf("skipping %s: %v", fn.Name(), err)
				}
				continue
			}
			res := sig.Results()
			err = checkType(res, res, true)
			if err != nil {
				if verbose {
					log.Printf("skipping %s: %v", fn.Name(), err)
				}
				continue
			}
			funcs = append(funcs, FuncInfo{
				Func:     fn,
				FuncDecl: fd,
			})

			walk(needUnpack, par, par)
			walk(needPack, res, res)
		}

	}

	// Check for mangled name collisions.
	seen := make(map[string]types.Type)
	for _, typ := range needUnpack {
		if typ2, ok := seen[Mangle(typ)]; ok {
			panic(fmt.Sprintf("mangled name collision in SEXP unwrapper: %s hits %s", typ, typ2))
		}
	}
	seen = make(map[string]types.Type)
	for _, typ := range needPack {
		if typ2, ok := seen[Mangle(typ)]; ok {
			panic(fmt.Sprintf("mangled name collision in SEXP builder: %s hits %s", typ, typ2))
		}
	}

	return &Info{Funcs: funcs, Unpackers: needUnpack, Packers: needPack}, nil
}

// TODO(kortschak): Handle recursive type definitions correctly.

func checkType(typ, named types.Type, warnRefs bool) error {
	switch typ := typ.(type) {
	case *types.Named:
		return checkType(typ.Underlying(), typ, warnRefs)

	case *types.Array:

	case *types.Basic:
		switch typ.Kind() {
		case types.Int64, types.Uint64:
			if typ == named {
				return fmt.Errorf("unhandled integer type %s", typ)
			}
			return fmt.Errorf("unhandled integer type %s (%s)", named, typ)
		}

	case *types.Chan:
		if typ == named {
			return fmt.Errorf("unhandled chan type %s", typ)
		}
		return fmt.Errorf("unhandled chan type %s (%s)", named, typ)

	case *types.Interface:
		if !IsError(named) {
			return fmt.Errorf("unhandled interface type %s", named)
		}

	case *types.Map:
		if !types.Identical(typ.Key().Underlying(), types.Typ[types.String]) {
			if typ == named {
				return fmt.Errorf("unhandled non-string keyed map type %s", typ)
			}
			return fmt.Errorf("unhandled non-string keyed map type %s (%s)", named, typ)
		}
		elem := typ.Elem()
		err := checkType(elem, elem, warnRefs)
		if err != nil {
			return err
		}

	case *types.Pointer:
		elem := typ.Elem()
		err := checkType(elem, elem, warnRefs)
		if err != nil {
			return err
		}
		if !warnRefs {
			break
		}
		if typ == named {
			log.Printf("warning: pointer type %s", typ)
		} else {
			log.Printf("warning: pointer type %s (%s)", named, typ)
		}

	case *types.Signature:
		if typ == named {
			return fmt.Errorf("unhandled function type with signature %s", typ)
		}
		return fmt.Errorf("unhandled function type with signature %s (%s)", named, typ)

	case *types.Slice:
		elem := typ.Elem()
		err := checkType(elem, elem, warnRefs)
		if err != nil {
			return err
		}
		if !warnRefs {
			break
		}
		if typ == named {
			log.Printf("warning: slice type %s", typ)
		} else {
			log.Printf("warning: slice type %s (%s)", named, typ)
		}

	case *types.Struct:
		for i := 0; i < typ.NumFields(); i++ {
			f := typ.Field(i).Type()
			err := checkType(f, f, warnRefs)
			if err != nil {
				return err
			}
		}

	case *types.Tuple:
		for i := 0; i < typ.Len(); i++ {
			f := typ.At(i).Type()
			err := checkType(f, f, warnRefs)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type unpackers map[string]types.Type

func (v unpackers) visit(typ types.Type) {
	if _, ok := typ.Underlying().(*types.Interface); ok {
		panic(fmt.Sprintf("unhandled input parameter type: %q", typ))
	}
	s := typ.String()
	if _, ok := v[s]; !ok {
		v[s] = typ
	}
	if also := v.also(typ); also != types.Invalid {
		v.visit(types.Typ[also])
	}
}

func (v unpackers) also(typ types.Type) types.BasicKind {
	basic, ok := typ.(*types.Basic)
	if !ok {
		return types.Invalid
	}
	// Make sure we have a complex128 value we can convert to complex64.
	if basic.Kind() == types.Complex64 {
		return types.Complex128
	}
	return types.Invalid
}

func (v unpackers) Types() []types.Type {
	typs := make([]types.Type, 0, len(v))
	for _, typ := range v {
		typs = append(typs, typ)
	}
	sort.Sort(byName(typs))
	return typs
}

type packers map[string]types.Type

func (v packers) visit(typ types.Type) {
	s := typ.String()
	if _, ok := v[s]; !ok {
		v[s] = typ
	}
}

func (v packers) NeedList() bool {
	for _, typ := range v {
		switch typ.(type) {
		case *types.Struct, *types.Map:
			return true
		}
	}
	return false
}

func (v packers) Types() []types.Type {
	typs := make([]types.Type, 0, len(v))
	for _, typ := range v {
		typs = append(typs, typ)
	}
	sort.Sort(byName(typs))
	return typs
}

type byName []types.Type

func (t byName) Len() int           { return len(t) }
func (t byName) Less(i, j int) bool { return Mangle(t[i]) < Mangle(t[j]) }
func (t byName) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

type visitor interface {
	visit(typ types.Type)
}

func walk(v visitor, typ, named types.Type) {
	switch typ := typ.(type) {
	case *types.Named:
		v.visit(typ)
		walk(v, typ.Underlying(), typ)

	case *types.Array:
		v.visit(typ)
		elem := typ.Elem()
		v.visit(types.NewSlice(elem)) // This will visit the element in the slice walk.

	case *types.Basic:
		switch typ.Kind() {
		case types.Int64, types.Uint64:
			if typ == named {
				panic(fmt.Sprintf("unhandled integer type %s", typ))
			}
			panic(fmt.Sprintf("unhandled integer type %s (%s)", named, typ))
		}
		v.visit(types.Typ[typ.Kind()])

	case *types.Chan:
		if typ == named {
			panic(fmt.Sprintf("unhandled chan type %s", typ))
		}
		panic(fmt.Sprintf("unhandled chan type %s (%s)", named, typ))

	case *types.Interface:
		if !IsError(named) {
			panic(fmt.Sprintf("unhandled interface type %s", named))
		}
		v.visit(types.Typ[types.String])
		v.visit(named)

	case *types.Map:
		if !types.Identical(typ.Key().Underlying(), types.Typ[types.String]) {
			if typ == named {
				panic(fmt.Sprintf("unhandled non-string keyed map type %s", typ))
			}
			panic(fmt.Sprintf("unhandled non-string keyed map type %s (%s)", named, typ))
		}
		// TODO(kortschak): De-alias the key and elem type in the map as well.
		v.visit(typ)
		key := typ.Key()
		walk(v, key, key)
		elem := typ.Elem()
		walk(v, elem, elem)

	case *types.Pointer:
		elem := typ.Elem()
		v.visit(typ)
		walk(v, elem, elem)

	case *types.Signature:
		if typ == named {
			panic(fmt.Sprintf("unhandled function type %s", typ))
		}
		panic(fmt.Sprintf("unhandled function type %s (%s)", named, typ))

	case *types.Slice:
		v.visit(typ)
		elem := typ.Elem()
		walk(v, elem, elem)

	case *types.Struct:
		v.visit(typ)
		for i := 0; i < typ.NumFields(); i++ {
			f := typ.Field(i).Type()
			walk(v, f, f)
		}

	case *types.Tuple:
		for i := 0; i < typ.Len(); i++ {
			f := typ.At(i).Type()
			walk(v, f, f)
		}
	}
}

func IsError(typ types.Type) bool {
	return types.Identical(typ, types.Universe.Lookup("error").Type())
}

func Mangle(typ types.Type) string {
	// FIXME(kortschak): This may lead to name collisions for complex unnamed types.
	runes := []rune(fmt.Sprintf("%T_%[1]s", typ))
	for i, r := range runes {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			runes[i] = '_'
		}
	}
	return string(runes)
}
