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

// packSEXPFuncGo returns the source of functions to pack the given Go-typed
// parameters into R SEXP values.
func packSEXPFuncGo(typs []types.Type) string {
	var buf bytes.Buffer
	for _, typ := range typs {
		fmt.Fprintf(&buf, "func packSEXP%s(p %s) C.SEXP {\n", pkg.Mangle(typ), nameOf(typ))
		packSEXPFuncBodyGo(&buf, typ)
		buf.WriteString("}\n\n")
	}
	return buf.String()
}

// packSEXPFuncGo returns the body of a function to pack the given Go-typed
// parameters into R SEXP values.
func packSEXPFuncBodyGo(buf *bytes.Buffer, typ types.Type) {
	switch typ := typ.(type) {
	case *types.Named:
		packNamed(buf, typ)

	case *types.Array:
		packArray(buf, typ)

	case *types.Basic:
		packBasic(buf, typ)

	case *types.Map:
		packMap(buf, typ)

	case *types.Pointer:
		packPointer(buf, typ)

	case *types.Slice:
		packSlice(buf, typ)

	case *types.Struct:
		packStruct(buf, typ)

	default:
		panic(fmt.Sprintf("unhandled type: %s", typ))
	}
}

func packNamed(buf *bytes.Buffer, typ *types.Named) {
	if pkg.IsError(typ) {
		fmt.Fprintf(buf, `	if p == nil {
		return C.R_NilValue
	}
	return packSEXP%s(p.Error())
`, pkg.Mangle(types.Typ[types.String]))
	} else {
		switch typ := typ.Underlying().(type) {
		case *types.Array, *types.Map, *types.Pointer, *types.Slice, *types.Struct:
			fmt.Fprintf(buf, "\treturn packSEXP%s(p)\n", pkg.Mangle(typ))
		default:
			fmt.Fprintf(buf, "\treturn packSEXP%s(%s(p))\n", pkg.Mangle(typ), typ)
		}
	}
}

func packArray(buf *bytes.Buffer, typ *types.Array) {
	fmt.Fprintf(buf, "\treturn packSEXP%s(p[:])\n", pkg.Mangle(types.NewSlice(typ.Elem())))
}

func packBasic(buf *bytes.Buffer, typ *types.Basic) {
	switch typ.Kind() {
	case types.Bool:
		fmt.Fprintf(buf, `	b := C.int(0)
	if p {
		b = 1
	}
	return C.ScalarLogical(b)
`)
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.Uint, types.Uint16, types.Uint32, types.Uint64:
		fmt.Fprintln(buf, "\treturn C.ScalarInteger(C.int(p))")
	case types.Uint8:
		fmt.Fprintln(buf, "\treturn C.ScalarRaw(C.Rbyte(p))")
	case types.Float64, types.Float32:
		fmt.Fprintln(buf, "\treturn C.ScalarReal(C.double(p))")
	case types.Complex128, types.Complex64:
		fmt.Fprintln(buf, "\treturn C.ScalarComplex(C.struct_Rcomplex{r: C.double(real(p)), i: C.double(imag(p))})")
	case types.String:
		fmt.Fprintln(buf, `	s := C.Rf_mkCharLenCE(C._GoStringPtr(p), C.int(len(p)), C.CE_UTF8)
	return C.ScalarString(s)`)
	default:
		panic(fmt.Sprintf("unhandled type: %s", typ))
	}
}

func packMap(buf *bytes.Buffer, typ *types.Map) {
	// TODO(kortschak): Handle named simple types properly.
	elem := typ.Elem()
	if basic, ok := elem.Underlying().(*types.Basic); ok {
		switch basic.Kind() {
		// TODO(kortschak): Make the fast path available
		// to []T where T is one of these kinds.
		case types.Int, types.Int8, types.Int16, types.Int32, types.Uint, types.Uint16, types.Uint32:
			// Maximum length array type for this element type.
			type a [1 << 47]int32
			fmt.Fprintf(buf, `	n := len(p)
	r := C.Rf_allocVector(C.%[1]s, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	s := (*[%[2]d]int32)(unsafe.Pointer(C.INTEGER(r)))[:len(p):len(p)]
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		s[i] = int32(v)
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
`, rTypeLabelFor(elem), len(&a{}))
			return

		case types.Uint8:
			// Maximum length array type for this element type.
			type a [1 << 49]byte
			fmt.Fprintf(buf, `	n := len(p)
	r := C.Rf_allocVector(C.%[1]s, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	s := (*[%[2]d]uint8)(unsafe.Pointer(C.RAW(r)))[:len(p):len(p)]
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		i++
	}
	copy(s, p)
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
`, rTypeLabelFor(elem), len(&a{}), pkg.Mangle(elem))
			return

		case types.Float32, types.Float64:
			// Maximum length array type for this element type.
			type a [1 << 46]float64
			fmt.Fprintf(buf, `	n := len(p)
	r := C.Rf_allocVector(C.%[1]s, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	s := (*[%[2]d]float64)(unsafe.Pointer(C.REAL(r)))[:len(p):len(p)]
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		s[i] = float64(v)
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
`, rTypeLabelFor(elem), len(&a{}))
			return

		case types.Complex64, types.Complex128:
			// Maximum length array type for this element type.
			type a [1 << 45]complex128
			fmt.Fprintf(buf, `	n := len(p)
	r := C.Rf_allocVector(C.%[1]s, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	s := (*[%[2]d]complex128)(unsafe.Pointer(C.COMPLEX(r)))[:len(p):len(p)]
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		s[i] = complex128(v)
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
`, rTypeLabelFor(elem), len(&a{}))
			return

		case types.String:
			fmt.Fprintf(buf, `	n := len(p)
	r := C.Rf_allocVector(C.%[1]s, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		C.SET_STRING_ELT(r, i, packSEXP%s(v))
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
`, rTypeLabelFor(elem), pkg.Mangle(elem))
			return

		case types.Bool:
			// Maximum length array type for this element type.
			type a [1 << 47]int32
			// FIXME(kortschak): Does Rf_allocVector return a
			// zeroed vector? If it does, the loop below doesn't
			// need the else clause.
			// Alternatively, convert the []bool to a []byte:
			//  for i, v := range *(*[]byte)(unsafe.Pointer(&p)) {
			//      s[i] = int32(v)
			//  }
			fmt.Fprintf(buf, `	n := len(p)
	r := C.Rf_allocVector(C.LGLSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	s := (*[%d]int32)(unsafe.Pointer(C.LOGICAL(r)))[:len(p):len(p)]
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		if v {
			s[i] = 1
		} else {
			s[i] = 0
		}
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
`, len(&a{}))
			return
		}
	}

	switch {
	case elem.String() == "error":
		fmt.Fprintf(buf, `	n := len(p)
	r := C.Rf_allocVector(C.%[1]s, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		if v == nil {
			C.SET_STRING_ELT(r, i, C.R_NilValue)
		} else {
			C.SET_STRING_ELT(r, i, packSEXP%[2]s(v))
		}
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
`, rTypeLabelFor(elem), pkg.Mangle(elem))

	default:
		if canBeNil(elem) {
			fmt.Fprintf(buf, `	n := len(p)
	r := C.Rf_allocVector(C.VECSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		if v == nil {
			C.SET_VECTOR_ELT(r, i, C.R_NilValue)
		} else {
			C.SET_VECTOR_ELT(r, i, packSEXP%s(v))
		}
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
`, pkg.Mangle(elem))
		} else {
			fmt.Fprintf(buf, `	n := len(p)
	r := C.Rf_allocVector(C.VECSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		C.SET_VECTOR_ELT(r, i, packSEXP%s(v))
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
`, pkg.Mangle(elem))
		}
	}
}

func packPointer(buf *bytes.Buffer, typ *types.Pointer) {
	fmt.Fprintf(buf, `	if p == nil {
		return C.R_NilValue
	}
	return packSEXP%s(*p)
`, pkg.Mangle(typ.Elem()))
}

func packSlice(buf *bytes.Buffer, typ *types.Slice) {
	// TODO(kortschak): Handle named simple types properly.
	elem := typ.Elem()
	if elem, ok := elem.(*types.Basic); ok {
		switch elem.Kind() {
		// TODO(kortschak): Make the fast path available
		// to []T where T is one of these kinds.
		case types.Int32, types.Uint32:
			// Maximum length array type for this element type.
			type a [1 << 47]int32
			fmt.Fprintf(buf, `	r := C.Rf_allocVector(C.INTSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[%d]%s)(unsafe.Pointer(C.INTEGER(r)))[:len(p)]
	copy(s, p)
	C.Rf_unprotect(1)
	return r
`, len(&a{}), nameOf(elem))
			return
		case types.Int, types.Int16, types.Uint, types.Uint16:
			// Maximum length array type for this element type.
			type a [1 << 47]int32
			fmt.Fprintf(buf, `	r := C.Rf_allocVector(C.INTSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[%[1]d]%[2]s)(unsafe.Pointer(C.INTEGER(r)))[:len(p)]
	for i, v := range p {
		s[i] = %[2]s(v)
	}
	C.Rf_unprotect(1)
	return r
`, len(&a{}), nameOf(types.Typ[types.Int32]))
			return
		case types.Int8, types.Uint8:
			// Maximum length array type for this element type.
			type a [1 << 49]byte
			fmt.Fprintf(buf, `	r := C.Rf_allocVector(C.RAWSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[%d]%s)(unsafe.Pointer(C.RAW(r)))[:len(p)]
	copy(s, p)
	C.Rf_unprotect(1)
	return r
`, len(&a{}), nameOf(elem))
			return
		case types.Float32:
			// Maximum length array type for this element type.
			type a [1 << 46]float64
			fmt.Fprintf(buf, `	r := C.Rf_allocVector(C.REALSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[%[1]d]%[2]s)(unsafe.Pointer(C.REAL(r)))[:len(p)]
	for i, v := range p {
		s[i] = %[2]s(v)
	}
	C.Rf_unprotect(1)
	return r
`, len(&a{}), nameOf(types.Typ[types.Float64]))
			return
		case types.Float64:
			// Maximum length array type for this element type.
			type a [1 << 46]float64
			fmt.Fprintf(buf, `	r := C.Rf_allocVector(C.REALSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[%d]%s)(unsafe.Pointer(C.REAL(r)))[:len(p)]
	copy(s, p)
	C.Rf_unprotect(1)
	return r
`, len(&a{}), nameOf(elem))
			return
		case types.Complex64:
			// Maximum length array type for this element type.
			type a [1 << 45]complex128
			fmt.Fprintf(buf, `	r := C.Rf_allocVector(C.CPLXSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[%[1]d]%[2]s)(unsafe.Pointer(C.CPLXSXP(r)))[:len(p)]
	for i, v := range p {
		s[i] = %[2]s(v)
	}
	C.Rf_unprotect(1)
	return r
`, len(&a{}), nameOf(types.Typ[types.Complex128]))
			return
		case types.Complex128:
			// Maximum length array type for this element type.
			type a [1 << 45]complex128
			fmt.Fprintf(buf, `	r := C.Rf_allocVector(C.CPLXSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[%d]%s)(unsafe.Pointer(C.CPLXSXP(r)))[:len(p)]
	copy(s, p)
	C.Rf_unprotect(1)
	return r
`, len(&a{}), nameOf(elem))
			return
		case types.Bool:
			// Maximum length array type for this element type.
			type a [1 << 47]int32
			// FIXME(kortschak): Does Rf_allocVector return a
			// zeroed vector? If it does, the loop below doesn't
			// need the else clause.
			// Alternatively, convert the []bool to a []byte:
			//  for i, v := range *(*[]byte)(unsafe.Pointer(&p)) {
			//      s[i] = int32(v)
			//  }
			fmt.Fprintf(buf, `	r := C.Rf_allocVector(C.LGLSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[%d]%s)(unsafe.Pointer(C.LOGICAL(r)))[:len(p)]
	for i, v := range p {
		if v {
			s[i] = 1
		} else {
			s[i] = 0
		}
	}
	C.Rf_unprotect(1)
	return r
`, len(&a{}), nameOf(elem))
			return
		case types.String:
			fmt.Fprint(buf, `	r := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	for i, v := range p {
		s := C.Rf_mkCharLenCE(C._GoStringPtr(string(v)), C.int(len(v)), C.CE_UTF8)
		C.SET_STRING_ELT(r, C.R_xlen_t(i), s)
	}
	C.Rf_unprotect(1)
	return r
`)
			return
		}
	}

	switch {
	case elem.String() == "error":
		fmt.Fprint(buf, `	r := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	for i, v := range p {
		s := C.R_NilValue
		if v != nil {
			s = C.Rf_mkCharLenCE(C._GoStringPtr(v), C.int(len(v)), C.CE_UTF8)
		}
		C.SET_STRING_ELT(r, C.R_xlen_t(i), s)
	}
	C.Rf_unprotect(1)
	return r
`)
	default:
		if canBeNil(elem) {
			fmt.Fprintf(buf, `	n := len(p)
	r := C.Rf_allocVector(C.VECSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	for i, v := range p {
		if v == nil {
			C.SET_VECTOR_ELT(r, C.R_xlen_t(i), C.R_NilValue)
		} else {
			C.SET_VECTOR_ELT(r, C.R_xlen_t(i), packSEXP%s(v))
		}
	}
	C.Rf_unprotect(1)
	return r
`, pkg.Mangle(elem))
		} else {
			fmt.Fprintf(buf, `	n := len(p)
	r := C.Rf_allocVector(C.VECSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	for i, v := range p {
		C.SET_VECTOR_ELT(r, C.R_xlen_t(i), packSEXP%s(v))
	}
	C.Rf_unprotect(1)
	return r
`, pkg.Mangle(elem))
		}
	}
}

func packStruct(buf *bytes.Buffer, typ *types.Struct) {
	n := typ.NumFields()
	fmt.Fprintf(buf, `	r := C.Rf_allocVector(C.VECSXP, %[1]d)
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, %[1]d)
	C.Rf_protect(names)
`, n)
	for i := 0; i < n; i++ {
		f := typ.Field(i)
		rName := targetFieldName(typ, i)
		elem := f.Type()
		fmt.Fprintf(buf, `	C.SET_STRING_ELT(names, %[1]d, C.Rf_mkCharLenCE(C._GoStringPtr("%[2]s"), %[3]d, C.CE_UTF8))
`, i, rName, len(rName))
		if canBeNil(elem) {
			fmt.Fprintf(buf, `	if v == nil {
		C.SET_VECTOR_ELT(r, %[1]d, C.R_NilValue)
	} else {
		C.SET_VECTOR_ELT(r, %[1]d, packSEXP%[2]s(p.%[3]s))
	}
`, i, pkg.Mangle(elem), f.Name())
		} else {
			fmt.Fprintf(buf, `	C.SET_VECTOR_ELT(r, %[1]d, packSEXP%[2]s(p.%[3]s))
`, i, pkg.Mangle(elem), f.Name())
		}
	}
	fmt.Fprintln(buf, `	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r`)
}

var typeLabelTable = map[string]string{
	"logical":   "LGLSXP",
	"integer":   "INTSXP",
	"double":    "REALSXP",
	"complex":   "CPLXSXP",
	"character": "STRSXP",
	"raw":       "RAWSXP",
	"list":      "VECSXP",
}

// rTypeLabelFor returns the R type label for the R atomic type
// corresponding to typ.
func rTypeLabelFor(typ types.Type) string {
	name, _, _ := rTypeOf(typ)
	label, ok := typeLabelTable[name]
	if !ok {
		return fmt.Sprintf("<%s>", typ)
	}
	return label
}

func canBeNil(typ types.Type) bool {
	switch typ.Underlying().(type) {
	case *types.Map, *types.Pointer, *types.Slice:
		return true
	default:
		return false
	}
}
