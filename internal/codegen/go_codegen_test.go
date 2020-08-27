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
	typ             types.Type
	wantUnpack      string
	wantUnpackNamed string
	wantPack        string
	wantPackNamed   string
}{
	// Basic types.
	{
		typ: types.Typ[types.String],
		wantUnpack: `func unpackSEXP_types_Basic_string(p C.SEXP) string {
	return C.R_gostring(p, 0)
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Basic_string(p))
}`,
		wantPack: `func packSEXP_types_Basic_string(p string) C.SEXP {
	s := C.Rf_mkCharLenCE(C._GoStringPtr(p), C.int(len(p)), C.CE_UTF8)
	return C.ScalarString(s)
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Basic_string(string(p))
}`,
	},

	{
		typ: types.Typ[types.Int32],
		wantUnpack: `func unpackSEXP_types_Basic_int32(p C.SEXP) int32 {
	return int32(*C.INTEGER(p))
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Basic_int32(p))
}`,
		wantPack: `func packSEXP_types_Basic_int32(p int32) C.SEXP {
	return C.ScalarInteger(C.int(p))
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Basic_int32(int32(p))
}`,
	},
	{
		typ: types.Universe.Lookup("rune").Type(),
		wantUnpack: `func unpackSEXP_types_Basic_rune(p C.SEXP) rune {
	return rune(*C.INTEGER(p))
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Basic_rune(p))
}`,
		wantPack: `func packSEXP_types_Basic_rune(p rune) C.SEXP {
	return C.ScalarInteger(C.int(p))
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Basic_rune(rune(p))
}`,
	},

	{
		typ: types.Typ[types.Uint8],
		wantUnpack: `func unpackSEXP_types_Basic_uint8(p C.SEXP) uint8 {
	return uint8(*C.RAW(p))
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Basic_uint8(p))
}`,
		wantPack: `func packSEXP_types_Basic_uint8(p uint8) C.SEXP {
	return C.ScalarRaw(C.Rbyte(p))
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Basic_uint8(uint8(p))
}`,
	},
	{
		typ: types.Universe.Lookup("byte").Type(),
		wantUnpack: `func unpackSEXP_types_Basic_byte(p C.SEXP) byte {
	return byte(*C.RAW(p))
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Basic_byte(p))
}`,
		wantPack: `func packSEXP_types_Basic_byte(p byte) C.SEXP {
	return C.ScalarRaw(C.Rbyte(p))
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Basic_byte(byte(p))
}`,
	},

	{
		typ: types.Typ[types.Float64],
		wantUnpack: `func unpackSEXP_types_Basic_float64(p C.SEXP) float64 {
	return float64(*C.REAL(p))
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Basic_float64(p))
}`,
		wantPack: `func packSEXP_types_Basic_float64(p float64) C.SEXP {
	return C.ScalarReal(C.double(p))
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Basic_float64(float64(p))
}`,
	},

	{
		typ: types.Typ[types.Complex128],
		wantUnpack: `func unpackSEXP_types_Basic_complex128(p C.SEXP) complex128 {
	return complex128(*(*complex128)(unsafe.Pointer(C.COMPLEX(p))))
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Basic_complex128(p))
}`,
		wantPack: `func packSEXP_types_Basic_complex128(p complex128) C.SEXP {
	return C.ScalarComplex(C.struct_Rcomplex{r: C.double(real(p)), i: C.double(imag(p))})
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Basic_complex128(complex128(p))
}`,
	},

	{
		typ: types.Typ[types.Bool],
		wantUnpack: `func unpackSEXP_types_Basic_bool(p C.SEXP) bool {
	return *C.RAW(p) == 1
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Basic_bool(p))
}`,
		wantPack: `func packSEXP_types_Basic_bool(p bool) C.SEXP {
	b := C.int(0)
	if p {
		b = 1
	}
	return C.ScalarLogical(b)
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Basic_bool(bool(p))
}`,
	},

	// Pointer types.
	{
		typ: types.NewPointer(types.Typ[types.String]),
		wantUnpack: `func unpackSEXP_types_Pointer__string(p C.SEXP) *string {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	r := unpackSEXP_types_Basic_string(p)
	return &r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Pointer__string(p))
}`,
		wantPack: `func packSEXP_types_Pointer__string(p *string) C.SEXP {
	if p == nil {
		return C.R_NilValue
	}
	return packSEXP_types_Basic_string(*p)
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Pointer__string((*string)(p))
}`,
	},

	{
		typ: types.NewPointer(types.Typ[types.Int32]),
		wantUnpack: `func unpackSEXP_types_Pointer__int32(p C.SEXP) *int32 {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	r := unpackSEXP_types_Basic_int32(p)
	return &r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Pointer__int32(p))
}`,
		wantPack: `func packSEXP_types_Pointer__int32(p *int32) C.SEXP {
	if p == nil {
		return C.R_NilValue
	}
	return packSEXP_types_Basic_int32(*p)
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Pointer__int32((*int32)(p))
}`,
	},
	{
		typ: types.NewPointer(types.Universe.Lookup("rune").Type()),
		wantUnpack: `func unpackSEXP_types_Pointer__rune(p C.SEXP) *rune {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	r := unpackSEXP_types_Basic_rune(p)
	return &r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Pointer__rune(p))
}`,
		wantPack: `func packSEXP_types_Pointer__rune(p *rune) C.SEXP {
	if p == nil {
		return C.R_NilValue
	}
	return packSEXP_types_Basic_rune(*p)
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Pointer__rune((*rune)(p))
}`,
	},

	{
		typ: types.NewPointer(types.Typ[types.Uint8]),
		wantUnpack: `func unpackSEXP_types_Pointer__uint8(p C.SEXP) *uint8 {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	r := unpackSEXP_types_Basic_uint8(p)
	return &r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Pointer__uint8(p))
}`,
		wantPack: `func packSEXP_types_Pointer__uint8(p *uint8) C.SEXP {
	if p == nil {
		return C.R_NilValue
	}
	return packSEXP_types_Basic_uint8(*p)
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Pointer__uint8((*uint8)(p))
}`,
	},
	{
		typ: types.NewPointer(types.Universe.Lookup("byte").Type()),
		wantUnpack: `func unpackSEXP_types_Pointer__byte(p C.SEXP) *byte {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	r := unpackSEXP_types_Basic_byte(p)
	return &r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Pointer__byte(p))
}`,
		wantPack: `func packSEXP_types_Pointer__byte(p *byte) C.SEXP {
	if p == nil {
		return C.R_NilValue
	}
	return packSEXP_types_Basic_byte(*p)
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Pointer__byte((*byte)(p))
}`,
	},

	{
		typ: types.NewPointer(types.Typ[types.Float64]),
		wantUnpack: `func unpackSEXP_types_Pointer__float64(p C.SEXP) *float64 {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	r := unpackSEXP_types_Basic_float64(p)
	return &r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Pointer__float64(p))
}`,
		wantPack: `func packSEXP_types_Pointer__float64(p *float64) C.SEXP {
	if p == nil {
		return C.R_NilValue
	}
	return packSEXP_types_Basic_float64(*p)
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Pointer__float64((*float64)(p))
}`,
	},

	{
		typ: types.NewPointer(types.Typ[types.Complex128]),
		wantUnpack: `func unpackSEXP_types_Pointer__complex128(p C.SEXP) *complex128 {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	r := unpackSEXP_types_Basic_complex128(p)
	return &r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Pointer__complex128(p))
}`,
		wantPack: `func packSEXP_types_Pointer__complex128(p *complex128) C.SEXP {
	if p == nil {
		return C.R_NilValue
	}
	return packSEXP_types_Basic_complex128(*p)
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Pointer__complex128((*complex128)(p))
}`,
	},

	{
		typ: types.NewPointer(types.Typ[types.Bool]),
		wantUnpack: `func unpackSEXP_types_Pointer__bool(p C.SEXP) *bool {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	r := unpackSEXP_types_Basic_bool(p)
	return &r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Pointer__bool(p))
}`,
		wantPack: `func packSEXP_types_Pointer__bool(p *bool) C.SEXP {
	if p == nil {
		return C.R_NilValue
	}
	return packSEXP_types_Basic_bool(*p)
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Pointer__bool((*bool)(p))
}`,
	},

	// Array types.
	{
		typ: types.NewArray(types.Typ[types.String], 10),
		wantUnpack: `func unpackSEXP_types_Array__10_string(p C.SEXP) [10]string {
	var a [10]string
	copy(a[:], unpackSEXP_types_Slice___string(p))
	return a
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Array__10_string(p))
}`,
		wantPack: `func packSEXP_types_Array__10_string(p [10]string) C.SEXP {
	return packSEXP_types_Slice___string(p[:])
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Array__10_string([10]string(p))
}`,
	},

	{
		typ: types.NewArray(types.Typ[types.Int32], 10),
		wantUnpack: `func unpackSEXP_types_Array__10_int32(p C.SEXP) [10]int32 {
	var a [10]int32
	copy(a[:], unpackSEXP_types_Slice___int32(p))
	return a
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Array__10_int32(p))
}`,
		wantPack: `func packSEXP_types_Array__10_int32(p [10]int32) C.SEXP {
	return packSEXP_types_Slice___int32(p[:])
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Array__10_int32([10]int32(p))
}`,
	},
	{
		typ: types.NewArray(types.Universe.Lookup("rune").Type(), 10),
		wantUnpack: `func unpackSEXP_types_Array__10_rune(p C.SEXP) [10]rune {
	var a [10]rune
	copy(a[:], unpackSEXP_types_Slice___rune(p))
	return a
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Array__10_rune(p))
}`,
		wantPack: `func packSEXP_types_Array__10_rune(p [10]rune) C.SEXP {
	return packSEXP_types_Slice___rune(p[:])
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Array__10_rune([10]rune(p))
}`,
	},

	{
		typ: types.NewArray(types.Typ[types.Uint8], 10),
		wantUnpack: `func unpackSEXP_types_Array__10_uint8(p C.SEXP) [10]uint8 {
	var a [10]uint8
	copy(a[:], unpackSEXP_types_Slice___uint8(p))
	return a
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Array__10_uint8(p))
}`,
		wantPack: `func packSEXP_types_Array__10_uint8(p [10]uint8) C.SEXP {
	return packSEXP_types_Slice___uint8(p[:])
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Array__10_uint8([10]uint8(p))
}`,
	},
	{
		typ: types.NewArray(types.Universe.Lookup("byte").Type(), 10),
		wantUnpack: `func unpackSEXP_types_Array__10_byte(p C.SEXP) [10]byte {
	var a [10]byte
	copy(a[:], unpackSEXP_types_Slice___byte(p))
	return a
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Array__10_byte(p))
}`,
		wantPack: `func packSEXP_types_Array__10_byte(p [10]byte) C.SEXP {
	return packSEXP_types_Slice___byte(p[:])
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Array__10_byte([10]byte(p))
}`,
	},

	{
		typ: types.NewArray(types.Typ[types.Float64], 10),
		wantUnpack: `func unpackSEXP_types_Array__10_float64(p C.SEXP) [10]float64 {
	var a [10]float64
	copy(a[:], unpackSEXP_types_Slice___float64(p))
	return a
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Array__10_float64(p))
}`,
		wantPack: `func packSEXP_types_Array__10_float64(p [10]float64) C.SEXP {
	return packSEXP_types_Slice___float64(p[:])
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Array__10_float64([10]float64(p))
}`,
	},

	{
		typ: types.NewArray(types.Typ[types.Complex128], 10),
		wantUnpack: `func unpackSEXP_types_Array__10_complex128(p C.SEXP) [10]complex128 {
	var a [10]complex128
	copy(a[:], unpackSEXP_types_Slice___complex128(p))
	return a
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Array__10_complex128(p))
}`,
		wantPack: `func packSEXP_types_Array__10_complex128(p [10]complex128) C.SEXP {
	return packSEXP_types_Slice___complex128(p[:])
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Array__10_complex128([10]complex128(p))
}`,
	},

	{
		typ: types.NewArray(types.Typ[types.Bool], 10),
		wantUnpack: `func unpackSEXP_types_Array__10_bool(p C.SEXP) [10]bool {
	var a [10]bool
	copy(a[:], unpackSEXP_types_Slice___bool(p))
	return a
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Array__10_bool(p))
}`,
		wantPack: `func packSEXP_types_Array__10_bool(p [10]bool) C.SEXP {
	return packSEXP_types_Slice___bool(p[:])
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Array__10_bool([10]bool(p))
}`,
	},

	// Slice types.
	{
		typ: types.NewSlice(types.Typ[types.String]),
		wantUnpack: `func unpackSEXP_types_Slice___string(p C.SEXP) []string {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := C.Rf_xlength(p)
	r := make([]string, n)
	for i := range r {
		r[i] = string(C.R_gostring(p, C.R_xlen_t(i)))
	}
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Slice___string(p))
}`,
		wantPack: `func packSEXP_types_Slice___string(p []string) C.SEXP {
	r := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	for i, v := range p {
		s := C.Rf_mkCharLenCE(C._GoStringPtr(string(v)), C.int(len(v)), C.CE_UTF8)
		C.SET_STRING_ELT(r, C.R_xlen_t(i), s)
	}
	C.Rf_unprotect(1)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Slice___string([]string(p))
}`,
	},

	{
		typ: types.NewSlice(types.Typ[types.Int32]),
		wantUnpack: `func unpackSEXP_types_Slice___int32(p C.SEXP) []int32 {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := C.Rf_xlength(p)
	return (*[140737488355328]int32)(unsafe.Pointer(C.INTEGER(p)))[:n:n]
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Slice___int32(p))
}`,
		wantPack: `func packSEXP_types_Slice___int32(p []int32) C.SEXP {
	r := C.Rf_allocVector(C.INTSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[140737488355328]int32)(unsafe.Pointer(C.INTEGER(r)))[:len(p):len(p)]
	copy(s, p)
	C.Rf_unprotect(1)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Slice___int32([]int32(p))
}`,
	},
	{
		typ: types.NewSlice(types.Universe.Lookup("rune").Type()),
		wantUnpack: `func unpackSEXP_types_Slice___rune(p C.SEXP) []rune {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := C.Rf_xlength(p)
	return (*[140737488355328]rune)(unsafe.Pointer(C.INTEGER(p)))[:n:n]
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Slice___rune(p))
}`,
		wantPack: `func packSEXP_types_Slice___rune(p []rune) C.SEXP {
	r := C.Rf_allocVector(C.INTSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[140737488355328]rune)(unsafe.Pointer(C.INTEGER(r)))[:len(p):len(p)]
	copy(s, p)
	C.Rf_unprotect(1)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Slice___rune([]rune(p))
}`,
	},

	{
		typ: types.NewSlice(types.Typ[types.Uint8]),
		wantUnpack: `func unpackSEXP_types_Slice___uint8(p C.SEXP) []uint8 {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := C.Rf_xlength(p)
	return (*[562949953421312]uint8)(unsafe.Pointer(C.RAW(p)))[:n:n]
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Slice___uint8(p))
}`,
		wantPack: `func packSEXP_types_Slice___uint8(p []uint8) C.SEXP {
	r := C.Rf_allocVector(C.RAWSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[562949953421312]uint8)(unsafe.Pointer(C.RAW(r)))[:len(p):len(p)]
	copy(s, p)
	C.Rf_unprotect(1)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Slice___uint8([]uint8(p))
}`,
	},
	{
		typ: types.NewSlice(types.Universe.Lookup("byte").Type()),
		wantUnpack: `func unpackSEXP_types_Slice___byte(p C.SEXP) []byte {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := C.Rf_xlength(p)
	return (*[562949953421312]byte)(unsafe.Pointer(C.RAW(p)))[:n:n]
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Slice___byte(p))
}`,
		wantPack: `func packSEXP_types_Slice___byte(p []byte) C.SEXP {
	r := C.Rf_allocVector(C.RAWSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[562949953421312]byte)(unsafe.Pointer(C.RAW(r)))[:len(p):len(p)]
	copy(s, p)
	C.Rf_unprotect(1)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Slice___byte([]byte(p))
}`,
	},

	{
		typ: types.NewSlice(types.Typ[types.Float64]),
		wantUnpack: `func unpackSEXP_types_Slice___float64(p C.SEXP) []float64 {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := C.Rf_xlength(p)
	return (*[70368744177664]float64)(unsafe.Pointer(C.REAL(p)))[:n:n]
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Slice___float64(p))
}`,
		wantPack: `func packSEXP_types_Slice___float64(p []float64) C.SEXP {
	r := C.Rf_allocVector(C.REALSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[70368744177664]float64)(unsafe.Pointer(C.REAL(r)))[:len(p):len(p)]
	copy(s, p)
	C.Rf_unprotect(1)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Slice___float64([]float64(p))
}`,
	},

	{
		typ: types.NewSlice(types.Typ[types.Complex128]),
		wantUnpack: `func unpackSEXP_types_Slice___complex128(p C.SEXP) []complex128 {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := C.Rf_xlength(p)
	return (*[35184372088832]complex128)(unsafe.Pointer(C.COMPLEX(p)))[:n:n]
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Slice___complex128(p))
}`,
		wantPack: `func packSEXP_types_Slice___complex128(p []complex128) C.SEXP {
	r := C.Rf_allocVector(C.CPLXSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[35184372088832]complex128)(unsafe.Pointer(C.CPLXSXP(r)))[:len(p):len(p)]
	copy(s, p)
	C.Rf_unprotect(1)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Slice___complex128([]complex128(p))
}`,
	},

	{
		typ: types.NewSlice(types.Typ[types.Bool]),
		wantUnpack: `func unpackSEXP_types_Slice___bool(p C.SEXP) []bool {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := C.Rf_xlength(p)
	r := make([]bool, n)
	for i, b := range (*[140737488355328]int32)(unsafe.Pointer(C.BOOL(p)))[:n] {
		r[i] = (b == 1)
	}
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Slice___bool(p))
}`,
		wantPack: `func packSEXP_types_Slice___bool(p []bool) C.SEXP {
	r := C.Rf_allocVector(C.LGLSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[140737488355328]bool)(unsafe.Pointer(C.LOGICAL(r)))[:len(p):len(p)]
	for i, v := range p {
		if v {
			s[i] = 1
		} else {
			s[i] = 0
		}
	}
	C.Rf_unprotect(1)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Slice___bool([]bool(p))
}`,
	},

	// Map types.
	{
		typ: types.NewMap(types.Typ[types.String], types.Typ[types.String]),
		wantUnpack: `func unpackSEXP_types_Map_map_string_string(p C.SEXP) map[string]string {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := int(C.Rf_xlength(p))
	r := make(map[string]string, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	for i := 0; i < n; i++ {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = string(C.R_gostring(p, C.R_xlen_t(i)))
	}
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Map_map_string_string(p))
}`,
		wantPack: `func packSEXP_types_Map_map_string_string(p map[string]string) C.SEXP {
	n := len(p)
	r := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		C.SET_STRING_ELT(r, i, packSEXP_types_Basic_string(v))
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Map_map_string_string(map[string]string(p))
}`,
	},

	{
		typ: types.NewMap(types.Typ[types.String], types.Typ[types.Int32]),
		wantUnpack: `func unpackSEXP_types_Map_map_string_int32(p C.SEXP) map[string]int32 {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := int(C.Rf_xlength(p))
	r := make(map[string]int32, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	values := (*[140737488355328]int32)(unsafe.Pointer(C.INTEGER(p)))[:n:n]
	for i, elem := range values {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = int32(elem)
	}
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Map_map_string_int32(p))
}`,
		wantPack: `func packSEXP_types_Map_map_string_int32(p map[string]int32) C.SEXP {
	n := len(p)
	r := C.Rf_allocVector(C.INTSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	s := (*[140737488355328]int32)(unsafe.Pointer(C.INTEGER(r)))[:len(p):len(p)]
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		s[i] = int32(v)
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Map_map_string_int32(map[string]int32(p))
}`,
	},
	{
		typ: types.NewMap(types.Typ[types.String], types.Universe.Lookup("rune").Type()),
		wantUnpack: `func unpackSEXP_types_Map_map_string_rune(p C.SEXP) map[string]rune {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := int(C.Rf_xlength(p))
	r := make(map[string]rune, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	values := (*[140737488355328]int32)(unsafe.Pointer(C.INTEGER(p)))[:n:n]
	for i, elem := range values {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = rune(elem)
	}
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Map_map_string_rune(p))
}`,
		wantPack: `func packSEXP_types_Map_map_string_rune(p map[string]rune) C.SEXP {
	n := len(p)
	r := C.Rf_allocVector(C.INTSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	s := (*[140737488355328]int32)(unsafe.Pointer(C.INTEGER(r)))[:len(p):len(p)]
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		s[i] = int32(v)
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Map_map_string_rune(map[string]rune(p))
}`,
	},

	{
		typ: types.NewMap(types.Typ[types.String], types.Typ[types.Uint8]),
		wantUnpack: `func unpackSEXP_types_Map_map_string_uint8(p C.SEXP) map[string]uint8 {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := int(C.Rf_xlength(p))
	r := make(map[string]uint8, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	values := (*[562949953421312]uint8)(unsafe.Pointer(C.RAW(p)))[:n:n]
	for i, elem := range values {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = elem
	}
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Map_map_string_uint8(p))
}`,
		wantPack: `func packSEXP_types_Map_map_string_uint8(p map[string]uint8) C.SEXP {
	n := len(p)
	r := C.Rf_allocVector(C.INTSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	s := (*[562949953421312]uint8)(unsafe.Pointer(C.RAW(r)))[:len(p):len(p)]
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		i++
	}
	copy(s, p)
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Map_map_string_uint8(map[string]uint8(p))
}`,
	},
	{
		typ: types.NewMap(types.Typ[types.String], types.Universe.Lookup("byte").Type()),
		wantUnpack: `func unpackSEXP_types_Map_map_string_byte(p C.SEXP) map[string]byte {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := int(C.Rf_xlength(p))
	r := make(map[string]byte, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	values := (*[562949953421312]byte)(unsafe.Pointer(C.RAW(p)))[:n:n]
	for i, elem := range values {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = elem
	}
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Map_map_string_byte(p))
}`,
		wantPack: `func packSEXP_types_Map_map_string_byte(p map[string]byte) C.SEXP {
	n := len(p)
	r := C.Rf_allocVector(C.INTSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	s := (*[562949953421312]uint8)(unsafe.Pointer(C.RAW(r)))[:len(p):len(p)]
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		i++
	}
	copy(s, p)
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Map_map_string_byte(map[string]byte(p))
}`,
	},

	{
		typ: types.NewMap(types.Typ[types.String], types.Typ[types.Float64]),
		wantUnpack: `func unpackSEXP_types_Map_map_string_float64(p C.SEXP) map[string]float64 {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := int(C.Rf_xlength(p))
	r := make(map[string]float64, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	values := (*[70368744177664]float64)(unsafe.Pointer(C.REAL(p)))[:n:n]
	for i, elem := range values {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = float64(elem)
	}
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Map_map_string_float64(p))
}`,
		wantPack: `func packSEXP_types_Map_map_string_float64(p map[string]float64) C.SEXP {
	n := len(p)
	r := C.Rf_allocVector(C.REALSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	s := (*[70368744177664]float64)(unsafe.Pointer(C.REAL(r)))[:len(p):len(p)]
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		s[i] = float64(v)
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Map_map_string_float64(map[string]float64(p))
}`,
	},

	{
		typ: types.NewMap(types.Typ[types.String], types.Typ[types.Complex128]),
		wantUnpack: `func unpackSEXP_types_Map_map_string_complex128(p C.SEXP) map[string]complex128 {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := int(C.Rf_xlength(p))
	r := make(map[string]complex128, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	values := (*[35184372088832]complex128)(unsafe.Pointer(C.COMPLEX(p)))[:n:n]
	for i, elem := range values {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = complex128(elem)
	}
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Map_map_string_complex128(p))
}`,
		wantPack: `func packSEXP_types_Map_map_string_complex128(p map[string]complex128) C.SEXP {
	n := len(p)
	r := C.Rf_allocVector(C.CPLXSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	s := (*[35184372088832]complex128)(unsafe.Pointer(C.COMPLEX(r)))[:len(p):len(p)]
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		s[i] = complex128(v)
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Map_map_string_complex128(map[string]complex128(p))
}`,
	},

	{
		typ: types.NewMap(types.Typ[types.String], types.Typ[types.Bool]),
		wantUnpack: `func unpackSEXP_types_Map_map_string_bool(p C.SEXP) map[string]bool {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := int(C.Rf_xlength(p))
	r := make(map[string]bool, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	values := (*[140737488355328]int32)(unsafe.Pointer(C.LOGICAL(p)))[:n:n]
	for i, elem := range values {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = (elem == 1)
	}
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Map_map_string_bool(p))
}`,
		wantPack: `func packSEXP_types_Map_map_string_bool(p map[string]bool) C.SEXP {
	n := len(p)
	r := C.Rf_allocVector(C.LGLSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	s := (*[140737488355328]int32)(unsafe.Pointer(C.LOGICAL(r)))[:len(p):len(p)]
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
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Map_map_string_bool(map[string]bool(p))
}`,
	},

	// Struct types.
	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Typ[types.String], false),
			types.NewField(0, mockPkg, "F2", types.Typ[types.String], false),
		}, []string{`rgo:"Rname"`}),
		wantUnpack: `func unpackSEXP_types_Struct_struct_F1_string__rgo___Rname_____F2_string_(p C.SEXP) struct{F1 string "rgo:\"Rname\""; F2 string} {
	switch n := C.Rf_xlength(p); {
	case n < 2:
		panic(` + "`missing list element for struct{F1 string \"rgo:\\\"Rname\\\"\"; F2 string}`" + `)
	case n > 2:
		err := C.CString(` + "`extra list element ignored for struct{F1 string \"rgo:\\\"Rname\\\"\"; F2 string}`" + `)
		C.R_error(err)
		C.free(unsafe.Pointer(err))
	}
	var r struct{F1 string "rgo:\"Rname\""; F2 string}
	var i C.int
	key_Rname := C.CString("Rname")
	defer C.free(unsafe.Pointer(key_Rname))
	i = C.getListElementIndex(p, key_Rname)
	if i < 0 {
		panic("no list element for field: F1")
	}
	r.F1 = unpackSEXP_types_Basic_string(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	key_F2 := C.CString("F2")
	defer C.free(unsafe.Pointer(key_F2))
	i = C.getListElementIndex(p, key_F2)
	if i < 0 {
		panic("no list element for field: F2")
	}
	r.F2 = unpackSEXP_types_Basic_string(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Struct_struct_F1_string__rgo___Rname_____F2_string_(p))
}`,
		wantPack: `func packSEXP_types_Struct_struct_F1_string__rgo___Rname_____F2_string_(p struct{F1 string "rgo:\"Rname\""; F2 string}) C.SEXP {
	r := C.allocList(2)
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, 2)
	C.Rf_protect(names)
	arg := r
	C.SET_STRING_ELT(names, 0, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`Rname`" + `), 5, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_string(p.F1))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 1, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`F2`" + `), 2, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_string(p.F2))
	C.setAttrib(r, packSEXP_types_Basic_string(` + "`names`" + `), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Struct_struct_F1_string__rgo___Rname_____F2_string_(struct{F1 string "rgo:\"Rname\""; F2 string}(p))
}`,
	},

	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Typ[types.Int32], false),
			types.NewField(0, mockPkg, "F2", types.Typ[types.Int32], false),
		}, []string{`rgo:"Rname"`}),
		wantUnpack: `func unpackSEXP_types_Struct_struct_F1_int32__rgo___Rname_____F2_int32_(p C.SEXP) struct{F1 int32 "rgo:\"Rname\""; F2 int32} {
	switch n := C.Rf_xlength(p); {
	case n < 2:
		panic(` + "`missing list element for struct{F1 int32 \"rgo:\\\"Rname\\\"\"; F2 int32}`" + `)
	case n > 2:
		err := C.CString(` + "`extra list element ignored for struct{F1 int32 \"rgo:\\\"Rname\\\"\"; F2 int32}`" + `)
		C.R_error(err)
		C.free(unsafe.Pointer(err))
	}
	var r struct{F1 int32 "rgo:\"Rname\""; F2 int32}
	var i C.int
	key_Rname := C.CString("Rname")
	defer C.free(unsafe.Pointer(key_Rname))
	i = C.getListElementIndex(p, key_Rname)
	if i < 0 {
		panic("no list element for field: F1")
	}
	r.F1 = unpackSEXP_types_Basic_int32(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	key_F2 := C.CString("F2")
	defer C.free(unsafe.Pointer(key_F2))
	i = C.getListElementIndex(p, key_F2)
	if i < 0 {
		panic("no list element for field: F2")
	}
	r.F2 = unpackSEXP_types_Basic_int32(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Struct_struct_F1_int32__rgo___Rname_____F2_int32_(p))
}`,
		wantPack: `func packSEXP_types_Struct_struct_F1_int32__rgo___Rname_____F2_int32_(p struct{F1 int32 "rgo:\"Rname\""; F2 int32}) C.SEXP {
	r := C.allocList(2)
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, 2)
	C.Rf_protect(names)
	arg := r
	C.SET_STRING_ELT(names, 0, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`Rname`" + `), 5, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_int32(p.F1))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 1, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`F2`" + `), 2, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_int32(p.F2))
	C.setAttrib(r, packSEXP_types_Basic_string(` + "`names`" + `), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Struct_struct_F1_int32__rgo___Rname_____F2_int32_(struct{F1 int32 "rgo:\"Rname\""; F2 int32}(p))
}`,
	},
	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Universe.Lookup("rune").Type(), false),
			types.NewField(0, mockPkg, "F2", types.Universe.Lookup("rune").Type(), false),
		}, []string{`rgo:"Rname"`}),
		wantUnpack: `func unpackSEXP_types_Struct_struct_F1_rune__rgo___Rname_____F2_rune_(p C.SEXP) struct{F1 rune "rgo:\"Rname\""; F2 rune} {
	switch n := C.Rf_xlength(p); {
	case n < 2:
		panic(` + "`missing list element for struct{F1 rune \"rgo:\\\"Rname\\\"\"; F2 rune}`" + `)
	case n > 2:
		err := C.CString(` + "`extra list element ignored for struct{F1 rune \"rgo:\\\"Rname\\\"\"; F2 rune}`" + `)
		C.R_error(err)
		C.free(unsafe.Pointer(err))
	}
	var r struct{F1 rune "rgo:\"Rname\""; F2 rune}
	var i C.int
	key_Rname := C.CString("Rname")
	defer C.free(unsafe.Pointer(key_Rname))
	i = C.getListElementIndex(p, key_Rname)
	if i < 0 {
		panic("no list element for field: F1")
	}
	r.F1 = unpackSEXP_types_Basic_rune(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	key_F2 := C.CString("F2")
	defer C.free(unsafe.Pointer(key_F2))
	i = C.getListElementIndex(p, key_F2)
	if i < 0 {
		panic("no list element for field: F2")
	}
	r.F2 = unpackSEXP_types_Basic_rune(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Struct_struct_F1_rune__rgo___Rname_____F2_rune_(p))
}`,
		wantPack: `func packSEXP_types_Struct_struct_F1_rune__rgo___Rname_____F2_rune_(p struct{F1 rune "rgo:\"Rname\""; F2 rune}) C.SEXP {
	r := C.allocList(2)
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, 2)
	C.Rf_protect(names)
	arg := r
	C.SET_STRING_ELT(names, 0, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`Rname`" + `), 5, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_rune(p.F1))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 1, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`F2`" + `), 2, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_rune(p.F2))
	C.setAttrib(r, packSEXP_types_Basic_string(` + "`names`" + `), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Struct_struct_F1_rune__rgo___Rname_____F2_rune_(struct{F1 rune "rgo:\"Rname\""; F2 rune}(p))
}`,
	},

	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Typ[types.Uint8], false),
			types.NewField(0, mockPkg, "F2", types.Typ[types.Uint8], false),
		}, []string{`rgo:"Rname"`}),
		wantUnpack: `func unpackSEXP_types_Struct_struct_F1_uint8__rgo___Rname_____F2_uint8_(p C.SEXP) struct{F1 uint8 "rgo:\"Rname\""; F2 uint8} {
	switch n := C.Rf_xlength(p); {
	case n < 2:
		panic(` + "`missing list element for struct{F1 uint8 \"rgo:\\\"Rname\\\"\"; F2 uint8}`" + `)
	case n > 2:
		err := C.CString(` + "`extra list element ignored for struct{F1 uint8 \"rgo:\\\"Rname\\\"\"; F2 uint8}`" + `)
		C.R_error(err)
		C.free(unsafe.Pointer(err))
	}
	var r struct{F1 uint8 "rgo:\"Rname\""; F2 uint8}
	var i C.int
	key_Rname := C.CString("Rname")
	defer C.free(unsafe.Pointer(key_Rname))
	i = C.getListElementIndex(p, key_Rname)
	if i < 0 {
		panic("no list element for field: F1")
	}
	r.F1 = unpackSEXP_types_Basic_uint8(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	key_F2 := C.CString("F2")
	defer C.free(unsafe.Pointer(key_F2))
	i = C.getListElementIndex(p, key_F2)
	if i < 0 {
		panic("no list element for field: F2")
	}
	r.F2 = unpackSEXP_types_Basic_uint8(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Struct_struct_F1_uint8__rgo___Rname_____F2_uint8_(p))
}`,
		wantPack: `func packSEXP_types_Struct_struct_F1_uint8__rgo___Rname_____F2_uint8_(p struct{F1 uint8 "rgo:\"Rname\""; F2 uint8}) C.SEXP {
	r := C.allocList(2)
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, 2)
	C.Rf_protect(names)
	arg := r
	C.SET_STRING_ELT(names, 0, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`Rname`" + `), 5, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_uint8(p.F1))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 1, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`F2`" + `), 2, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_uint8(p.F2))
	C.setAttrib(r, packSEXP_types_Basic_string(` + "`names`" + `), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Struct_struct_F1_uint8__rgo___Rname_____F2_uint8_(struct{F1 uint8 "rgo:\"Rname\""; F2 uint8}(p))
}`,
	},
	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Universe.Lookup("byte").Type(), false),
			types.NewField(0, mockPkg, "F2", types.Universe.Lookup("byte").Type(), false),
		}, []string{`rgo:"Rname"`}),
		wantUnpack: `func unpackSEXP_types_Struct_struct_F1_byte__rgo___Rname_____F2_byte_(p C.SEXP) struct{F1 byte "rgo:\"Rname\""; F2 byte} {
	switch n := C.Rf_xlength(p); {
	case n < 2:
		panic(` + "`missing list element for struct{F1 byte \"rgo:\\\"Rname\\\"\"; F2 byte}`" + `)
	case n > 2:
		err := C.CString(` + "`extra list element ignored for struct{F1 byte \"rgo:\\\"Rname\\\"\"; F2 byte}`" + `)
		C.R_error(err)
		C.free(unsafe.Pointer(err))
	}
	var r struct{F1 byte "rgo:\"Rname\""; F2 byte}
	var i C.int
	key_Rname := C.CString("Rname")
	defer C.free(unsafe.Pointer(key_Rname))
	i = C.getListElementIndex(p, key_Rname)
	if i < 0 {
		panic("no list element for field: F1")
	}
	r.F1 = unpackSEXP_types_Basic_byte(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	key_F2 := C.CString("F2")
	defer C.free(unsafe.Pointer(key_F2))
	i = C.getListElementIndex(p, key_F2)
	if i < 0 {
		panic("no list element for field: F2")
	}
	r.F2 = unpackSEXP_types_Basic_byte(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Struct_struct_F1_byte__rgo___Rname_____F2_byte_(p))
}`,
		wantPack: `func packSEXP_types_Struct_struct_F1_byte__rgo___Rname_____F2_byte_(p struct{F1 byte "rgo:\"Rname\""; F2 byte}) C.SEXP {
	r := C.allocList(2)
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, 2)
	C.Rf_protect(names)
	arg := r
	C.SET_STRING_ELT(names, 0, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`Rname`" + `), 5, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_byte(p.F1))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 1, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`F2`" + `), 2, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_byte(p.F2))
	C.setAttrib(r, packSEXP_types_Basic_string(` + "`names`" + `), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Struct_struct_F1_byte__rgo___Rname_____F2_byte_(struct{F1 byte "rgo:\"Rname\""; F2 byte}(p))
}`,
	},

	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Typ[types.Float64], false),
			types.NewField(0, mockPkg, "F2", types.Typ[types.Float64], false),
		}, []string{`rgo:"Rname"`}),
		wantUnpack: `func unpackSEXP_types_Struct_struct_F1_float64__rgo___Rname_____F2_float64_(p C.SEXP) struct{F1 float64 "rgo:\"Rname\""; F2 float64} {
	switch n := C.Rf_xlength(p); {
	case n < 2:
		panic(` + "`missing list element for struct{F1 float64 \"rgo:\\\"Rname\\\"\"; F2 float64}`" + `)
	case n > 2:
		err := C.CString(` + "`extra list element ignored for struct{F1 float64 \"rgo:\\\"Rname\\\"\"; F2 float64}`" + `)
		C.R_error(err)
		C.free(unsafe.Pointer(err))
	}
	var r struct{F1 float64 "rgo:\"Rname\""; F2 float64}
	var i C.int
	key_Rname := C.CString("Rname")
	defer C.free(unsafe.Pointer(key_Rname))
	i = C.getListElementIndex(p, key_Rname)
	if i < 0 {
		panic("no list element for field: F1")
	}
	r.F1 = unpackSEXP_types_Basic_float64(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	key_F2 := C.CString("F2")
	defer C.free(unsafe.Pointer(key_F2))
	i = C.getListElementIndex(p, key_F2)
	if i < 0 {
		panic("no list element for field: F2")
	}
	r.F2 = unpackSEXP_types_Basic_float64(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Struct_struct_F1_float64__rgo___Rname_____F2_float64_(p))
}`,
		wantPack: `func packSEXP_types_Struct_struct_F1_float64__rgo___Rname_____F2_float64_(p struct{F1 float64 "rgo:\"Rname\""; F2 float64}) C.SEXP {
	r := C.allocList(2)
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, 2)
	C.Rf_protect(names)
	arg := r
	C.SET_STRING_ELT(names, 0, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`Rname`" + `), 5, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_float64(p.F1))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 1, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`F2`" + `), 2, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_float64(p.F2))
	C.setAttrib(r, packSEXP_types_Basic_string(` + "`names`" + `), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Struct_struct_F1_float64__rgo___Rname_____F2_float64_(struct{F1 float64 "rgo:\"Rname\""; F2 float64}(p))
}`,
	},

	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Typ[types.Complex128], false),
			types.NewField(0, mockPkg, "F2", types.Typ[types.Complex128], false),
		}, []string{`rgo:"Rname"`}),
		wantUnpack: `func unpackSEXP_types_Struct_struct_F1_complex128__rgo___Rname_____F2_complex128_(p C.SEXP) struct{F1 complex128 "rgo:\"Rname\""; F2 complex128} {
	switch n := C.Rf_xlength(p); {
	case n < 2:
		panic(` + "`missing list element for struct{F1 complex128 \"rgo:\\\"Rname\\\"\"; F2 complex128}`" + `)
	case n > 2:
		err := C.CString(` + "`extra list element ignored for struct{F1 complex128 \"rgo:\\\"Rname\\\"\"; F2 complex128}`" + `)
		C.R_error(err)
		C.free(unsafe.Pointer(err))
	}
	var r struct{F1 complex128 "rgo:\"Rname\""; F2 complex128}
	var i C.int
	key_Rname := C.CString("Rname")
	defer C.free(unsafe.Pointer(key_Rname))
	i = C.getListElementIndex(p, key_Rname)
	if i < 0 {
		panic("no list element for field: F1")
	}
	r.F1 = unpackSEXP_types_Basic_complex128(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	key_F2 := C.CString("F2")
	defer C.free(unsafe.Pointer(key_F2))
	i = C.getListElementIndex(p, key_F2)
	if i < 0 {
		panic("no list element for field: F2")
	}
	r.F2 = unpackSEXP_types_Basic_complex128(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Struct_struct_F1_complex128__rgo___Rname_____F2_complex128_(p))
}`,
		wantPack: `func packSEXP_types_Struct_struct_F1_complex128__rgo___Rname_____F2_complex128_(p struct{F1 complex128 "rgo:\"Rname\""; F2 complex128}) C.SEXP {
	r := C.allocList(2)
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, 2)
	C.Rf_protect(names)
	arg := r
	C.SET_STRING_ELT(names, 0, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`Rname`" + `), 5, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_complex128(p.F1))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 1, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`F2`" + `), 2, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_complex128(p.F2))
	C.setAttrib(r, packSEXP_types_Basic_string(` + "`names`" + `), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Struct_struct_F1_complex128__rgo___Rname_____F2_complex128_(struct{F1 complex128 "rgo:\"Rname\""; F2 complex128}(p))
}`,
	},

	{
		typ: types.NewStruct([]*types.Var{
			types.NewField(0, mockPkg, "F1", types.Typ[types.Bool], false),
			types.NewField(0, mockPkg, "F2", types.Typ[types.Bool], false),
		}, []string{`rgo:"Rname"`}),
		wantUnpack: `func unpackSEXP_types_Struct_struct_F1_bool__rgo___Rname_____F2_bool_(p C.SEXP) struct{F1 bool "rgo:\"Rname\""; F2 bool} {
	switch n := C.Rf_xlength(p); {
	case n < 2:
		panic(` + "`missing list element for struct{F1 bool \"rgo:\\\"Rname\\\"\"; F2 bool}`" + `)
	case n > 2:
		err := C.CString(` + "`extra list element ignored for struct{F1 bool \"rgo:\\\"Rname\\\"\"; F2 bool}`" + `)
		C.R_error(err)
		C.free(unsafe.Pointer(err))
	}
	var r struct{F1 bool "rgo:\"Rname\""; F2 bool}
	var i C.int
	key_Rname := C.CString("Rname")
	defer C.free(unsafe.Pointer(key_Rname))
	i = C.getListElementIndex(p, key_Rname)
	if i < 0 {
		panic("no list element for field: F1")
	}
	r.F1 = unpackSEXP_types_Basic_bool(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	key_F2 := C.CString("F2")
	defer C.free(unsafe.Pointer(key_F2))
	i = C.getListElementIndex(p, key_F2)
	if i < 0 {
		panic("no list element for field: F2")
	}
	r.F2 = unpackSEXP_types_Basic_bool(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	return r
}`,
		wantUnpackNamed: `func unpackSEXP_types_Named_path_to_pkg_T(p C.SEXP) pkg.T {
	return pkg.T(unpackSEXP_types_Struct_struct_F1_bool__rgo___Rname_____F2_bool_(p))
}`,
		wantPack: `func packSEXP_types_Struct_struct_F1_bool__rgo___Rname_____F2_bool_(p struct{F1 bool "rgo:\"Rname\""; F2 bool}) C.SEXP {
	r := C.allocList(2)
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, 2)
	C.Rf_protect(names)
	arg := r
	C.SET_STRING_ELT(names, 0, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`Rname`" + `), 5, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_bool(p.F1))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 1, C.Rf_mkCharLenCE(C._GoStringPtr(` + "`F2`" + `), 2, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_bool(p.F2))
	C.setAttrib(r, packSEXP_types_Basic_string(` + "`names`" + `), names)
	C.Rf_unprotect(2)
	return r
}`,
		wantPackNamed: `func packSEXP_types_Named_path_to_pkg_T(p pkg.T) C.SEXP {
	return packSEXP_types_Struct_struct_F1_bool__rgo___Rname_____F2_bool_(struct{F1 bool "rgo:\"Rname\""; F2 bool}(p))
}`,
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
		for j, typ := range typs {
			got := []byte(strings.TrimSpace(unpackSEXPFuncGo([]types.Type{typ})))
			if (j == 0 && string(got) != test.wantUnpack) || (j == 1 && string(got) != test.wantUnpackNamed) {
				t.Errorf("unexpected result for test %d %s:\ngot:\n%s\nwant:\n%s", i, typ, got, test.wantUnpack)
			}

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
		for j, typ := range typs {
			got := []byte(strings.TrimSpace(packSEXPFuncGo([]types.Type{typ})))
			if (j == 0 && string(got) != test.wantPack) || (j == 1 && string(got) != test.wantPackNamed) {
				t.Errorf("unexpected result for test %d %s:\ngot:\n%s\nwant:\n%s", i, typ, got, test.wantUnpack)
			}

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
