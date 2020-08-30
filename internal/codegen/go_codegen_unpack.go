// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codegen

import (
	"bytes"
	"fmt"
	"go/types"

	"github.com/rgonomic/rgo/internal/pkg"
)

// unpackSEXPFuncGo returns the source of functions to unpack R SEXP parameters
// into the given Go types.
func unpackSEXPFuncGo(typs []types.Type) string {
	var buf bytes.Buffer
	for _, typ := range typs {
		fmt.Fprintf(&buf, "func unpackSEXP%s(p C.SEXP) %s {\n", pkg.Mangle(typ), nameOf(typ))
		unpackSEXPFuncBodyGo(&buf, typ)
		buf.WriteString("}\n\n")
	}
	return buf.String()
}

// unpackSEXPFuncBodyGo returns the body of a function to unpack R SEXP parameters
// into the given Go types.
func unpackSEXPFuncBodyGo(buf *bytes.Buffer, typ types.Type) {
	switch typ := typ.(type) {
	case *types.Named:
		unpackNamed(buf, typ)

	case *types.Array:
		unpackArray(buf, typ)

	case *types.Basic:
		unpackBasic(buf, typ)

	case *types.Map:
		unpackMap(buf, typ)

	case *types.Pointer:
		unpackPointer(buf, typ)

	case *types.Slice:
		unpackSlice(buf, typ)

	case *types.Struct:
		unpackStruct(buf, typ)

	default:
		panic(fmt.Sprintf("unhandled type: %s", typ))
	}
}

func unpackNamed(buf *bytes.Buffer, typ *types.Named) {
	switch under := typ.Underlying().(type) {
	case *types.Array, *types.Map, *types.Pointer, *types.Slice, *types.Struct:
		fmt.Fprintf(buf, "\treturn unpackSEXP%s(p)\n", pkg.Mangle(under))
	default:
		fmt.Fprintf(buf, "\treturn %s(unpackSEXP%s(p))\n", nameOf(typ), pkg.Mangle(under))
	}
}

func unpackArray(buf *bytes.Buffer, typ *types.Array) {
	// TODO(kortschak): Only do this for [n]int32, [n]float64, [n]complex128 and [n]byte.
	// Otherwise we have a double copy.
	fmt.Fprintf(buf, `	var a %s
	copy(a[:], unpackSEXP%s(p))
	return a
`, typ, pkg.Mangle(types.NewSlice(typ.Elem())))
}

func unpackBasic(buf *bytes.Buffer, typ *types.Basic) {
	switch typ.Kind() {
	case types.Bool:
		fmt.Fprintln(buf, "\treturn *C.RAW(p) == 1")
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.Uint, types.Uint16, types.Uint32, types.Uint64:
		fmt.Fprintf(buf, "\treturn %s(*C.INTEGER(p))\n", nameOf(typ))
	case types.Uint8:
		fmt.Fprintf(buf, "\treturn %s(*C.RAW(p))\n", nameOf(typ))
	case types.Float64, types.Float32:
		fmt.Fprintf(buf, "\treturn %s(*C.REAL(p))\n", nameOf(typ))
	case types.Complex128:
		fmt.Fprintf(buf, "\treturn %s(*(*complex128)(unsafe.Pointer(C.COMPLEX(p))))\n", nameOf(typ))
	case types.Complex64:
		fmt.Fprintf(buf, "\treturn %s(unpackSEXP%s(p))\n", nameOf(typ), pkg.Mangle(types.Typ[types.Complex128]))
	case types.String:
		fmt.Fprintln(buf, "\treturn C.R_gostring(p, 0)")
	case types.UnsafePointer:
		fmt.Fprintln(buf, "\treturn unsafe.Pointer(p)")
	default:
		panic(fmt.Sprintf("unhandled type: %s", typ))
	}
}

func unpackMap(buf *bytes.Buffer, typ *types.Map) {
	fmt.Fprintf(buf, `	if C.Rf_isNull(p) != 0 {
		return nil
	}
`)
	elem := typ.Elem()
	if basic, ok := elem.Underlying().(*types.Basic); ok {
		switch basic.Kind() {
		// TODO(kortschak): Make the fast path available
		// to []T where T is one of these kinds.
		case types.Int, types.Int8, types.Int16, types.Int32, types.Uint, types.Uint16, types.Uint32:
			// Maximum length array type for this element type.
			type a [1 << 47]int32
			fmt.Fprintf(buf, `	n := int(C.Rf_xlength(p))
	r := make(map[string]%[2]s, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	if names == C.R_NilValue {
		panic("no names attribute for map keys")
	}
	values := (*[%[1]d]int32)(unsafe.Pointer(C.INTEGER(p)))[:n:n]
	for i, elem := range values {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = %[2]s(elem)
	}
	return r
`, len(&a{}), nameOf(elem))
			return
		case types.Uint8:
			// Maximum length array type for this element type.
			type a [1 << 49]byte
			fmt.Fprintf(buf, `	n := int(C.Rf_xlength(p))
	r := make(map[string]%[2]s, n)
	if names == C.R_NilValue {
		panic("no names attribute for map keys")
	}
	names := C.getAttrib(p, C.R_NamesSymbol)
	values := (*[%[1]d]%[2]s)(unsafe.Pointer(C.RAW(p)))[:n:n]
	for i, elem := range values {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = elem
	}
	return r
`, len(&a{}), nameOf(elem))
			return
		case types.Float32, types.Float64:
			// Maximum length array type for this element type.
			type a [1 << 46]float64
			fmt.Fprintf(buf, `	n := int(C.Rf_xlength(p))
	r := make(map[string]%[2]s, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	if names == C.R_NilValue {
		panic("no names attribute for map keys")
	}
	values := (*[%[1]d]float64)(unsafe.Pointer(C.REAL(p)))[:n:n]
	for i, elem := range values {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = %[2]s(elem)
	}
	return r
`, len(&a{}), nameOf(elem))
			return
		case types.Complex64, types.Complex128:
			// Maximum length array type for this element type.
			type a [1 << 45]complex128
			fmt.Fprintf(buf, `	n := int(C.Rf_xlength(p))
	r := make(map[string]%[2]s, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	if names == C.R_NilValue {
		panic("no names attribute for map keys")
	}
	values := (*[%[1]d]complex128)(unsafe.Pointer(C.COMPLEX(p)))[:n:n]
	for i, elem := range values {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = %[2]s(elem)
	}
	return r
`, len(&a{}), nameOf(elem))
			return
		case types.Bool:
			// Maximum length array type for this element type.
			type a [1 << 47]int32
			fmt.Fprintf(buf, `	n := int(C.Rf_xlength(p))
	r := make(map[string]%[2]s, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	if names == C.R_NilValue {
		panic("no names attribute for map keys")
	}
	values := (*[%[1]d]int32)(unsafe.Pointer(C.LOGICAL(p)))[:n:n]
	for i, elem := range values {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = (elem == 1)
	}
	return r
`, len(&a{}), nameOf(elem))
			return
		case types.String:
			fmt.Fprintf(buf, `	n := int(C.Rf_xlength(p))
	r := make(map[string]%[1]s, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	if names == C.R_NilValue {
		panic("no names attribute for map keys")
	}
	for i := 0; i < n; i++ {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = %[1]s(C.R_gostring(p, C.R_xlen_t(i)))
	}
	return r
`, nameOf(elem))
			return
		}
	}

	fmt.Fprintf(buf, `	n := int(C.Rf_xlength(p))
	r := make(map[string]%s, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	if names == C.R_NilValue {
		panic("no names attribute for map keys")
	}
	for i := 0; i < n; i++ {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = unpackSEXP%s(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	}
	return r
`, nameOf(elem), pkg.Mangle(elem))
}

func unpackPointer(buf *bytes.Buffer, typ *types.Pointer) {
	fmt.Fprintf(buf, `	if C.Rf_isNull(p) != 0 {
		return nil
	}
	r := unpackSEXP%s(p)
	return &r
`, pkg.Mangle(typ.Elem()))
}

func unpackSlice(buf *bytes.Buffer, typ *types.Slice) {
	// TODO(kortschak): Use unsafe.Slice when it exists.

	fmt.Fprintf(buf, `	if C.Rf_isNull(p) != 0 {
		return nil
	}
`)
	elem := typ.Elem()
	if elem, ok := elem.(*types.Basic); ok {
		switch elem.Kind() {
		// TODO(kortschak): Make the fast path available
		// to []T where T is one of these kinds.
		case types.Int32, types.Uint32:
			// Maximum length array type for this element type.
			type a [1 << 47]int32
			fmt.Fprintf(buf, `	n := C.Rf_xlength(p)
	return (*[%d]%s)(unsafe.Pointer(C.INTEGER(p)))[:n]
`, len(&a{}), nameOf(elem))
			return
		case types.Int, types.Int16, types.Uint, types.Uint16:
			// Maximum length array type for this element type.
			type a [1 << 47]int32
			fmt.Fprintf(buf, `	n := C.Rf_xlength(p)
	r := make(%s, n)
	for i, v := range (*[%d]%s)(unsafe.Pointer(C.INTEGER(p)))[:n] {
		r[i] = %s(v)
	}
	return r
`, nameOf(typ), len(&a{}), nameOf(types.Typ[types.Int32]), elem)
			return
		case types.Int8, types.Uint8:
			// Maximum length array type for this element type.
			type a [1 << 49]byte
			fmt.Fprintf(buf, `	n := C.Rf_xlength(p)
	return (*[%d]%s)(unsafe.Pointer(C.RAW(p)))[:n]
`, len(&a{}), nameOf(elem))
			return
		case types.Float32:
			// Maximum length array type for this element type.
			type a [1 << 46]float64
			fmt.Fprintf(buf, `	n := C.Rf_xlength(p)
	r := make(%s, n)
	for i, v := range (*[%d]%s)(unsafe.Pointer(C.REAL(p)))[:n] {
		r[i] = %s(v)
	}
	return r
`, nameOf(typ), len(&a{}), nameOf(types.Typ[types.Float64]), elem)
			return
		case types.Float64:
			// Maximum length array type for this element type.
			type a [1 << 46]float64
			fmt.Fprintf(buf, `	n := C.Rf_xlength(p)
	return (*[%d]%s)(unsafe.Pointer(C.REAL(p)))[:n]
`, len(&a{}), nameOf(elem))
			return
		case types.Complex64:
			// Maximum length array type for this element type.
			type a [1 << 45]complex128
			fmt.Fprintf(buf, `	n := C.Rf_xlength(p)
	r := make(%s, n)
	for i, v := range (*[%d]%s)(unsafe.Pointer(C.COMPLEX(p)))[:n] {
		r[i] = %s(v)
	}
	return r
`, nameOf(typ), len(&a{}), nameOf(types.Typ[types.Complex128]), elem)
			return
		case types.Complex128:
			// Maximum length array type for this element type.
			type a [1 << 45]complex128
			fmt.Fprintf(buf, `	n := C.Rf_xlength(p)
	return (*[%d]%s)(unsafe.Pointer(C.COMPLEX(p)))[:n]
`, len(&a{}), nameOf(elem))
			return
		case types.Bool:
			// Maximum length array type for this element type.
			type a [1 << 47]int32
			fmt.Fprintf(buf, `	n := C.Rf_xlength(p)
	r := make(%s, n)
	for i, b := range (*[%d]%s)(unsafe.Pointer(C.BOOL(p)))[:n] {
		r[i] = (b == 1)
	}
	return r
`, nameOf(typ), len(&a{}), nameOf(types.Typ[types.Int32]))
			return
		case types.String:
			fmt.Fprintf(buf, `	n := C.Rf_xlength(p)
	r := make(%s, n)
	for i := range r {
		r[i] = %s(C.R_gostring(p, C.R_xlen_t(i)))
	}
	return r
`, nameOf(typ), nameOf(elem))
			return
		}
	}
	fmt.Fprintf(buf, `	n := C.Rf_xlength(p)
	r := make(%s, n)
	for i := range r {
		r[i] = unpackSEXP%s(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	}
	return r
`, nameOf(typ), pkg.Mangle(elem))
}

func unpackStruct(buf *bytes.Buffer, typ *types.Struct) {
	n := typ.NumFields()
	fmt.Fprintf(buf, `	switch n := C.Rf_xlength(p); {
	case n < %[1]d:
		panic(`+"`missing list element for %[2]s`"+`)
	case n > %[1]d:
		err := C.CString(`+"`extra list element ignored for %[2]s`"+`)
		C.R_error(err)
		C.free(unsafe.Pointer(err))
	}
	var r %[2]s
	var i C.int
`, n, nameOf(typ))
	for i := 0; i < n; i++ {
		f := typ.Field(i)

		fmt.Fprintf(buf, `	key_%s := C.CString("%[1]s")
	defer C.free(unsafe.Pointer(key_%[1]s))
	i = C.getListElementIndex(p, key_%[1]s)
	if i < 0 {
		panic("no list element name for field: %[2]s")
	}
	r.%[2]s = unpackSEXP%s(C.VECTOR_ELT(p, C.R_xlen_t(i)))
`, targetFieldName(typ, i), f.Name(), pkg.Mangle(f.Type()))
	}
	fmt.Fprintln(buf, "\treturn r")
}
