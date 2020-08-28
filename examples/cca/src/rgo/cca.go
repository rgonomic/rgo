// Code generated by rgnonomic/rgo; DO NOT EDIT.

package main

/*
#define USE_RINTERNALS
#include <R.h>
#include <Rinternals.h>
extern void R_error(char *s);

// TODO(kortschak): Only emit these when needed.
extern Rboolean Rf_isNull(SEXP s);
extern _GoString_ R_gostring(SEXP x, R_xlen_t i);
extern int getListElementIndex(SEXP list, const char *str);
*/
import "C"

import (
	"fmt"
	"unsafe"

	"gonum.org/v1/gonum/blas/blas64"

	"github.com/rgonomic/rgo/examples/cca"
)

//export Wrapped_CCA
func Wrapped_CCA(_R_x, _R_y, _R_weights C.SEXP) C.SEXP {
	defer func() {
		r := recover()
		if r != nil {
			err := C.CString(fmt.Sprint(r))
			C.R_error(err)
			C.free(unsafe.Pointer(err))
		}
	}()

	_p0 := unpackSEXP_types_Named_gonum_org_v1_gonum_blas_blas64_GeneralCols(_R_x)
	_p1 := unpackSEXP_types_Named_gonum_org_v1_gonum_blas_blas64_GeneralCols(_R_y)
	_p2 := unpackSEXP_types_Slice___float64(_R_weights)
	_r0, _r1, _r2, _r3, _r4, _r5 := cca.CCA(_p0, _p1, _p2)
	return packSEXP_CCA(_r0, _r1, _r2, _r3, _r4, _r5)
}

func packSEXP_CCA(ccors []float64, pVecs blas64.GeneralCols, qVecs blas64.GeneralCols, phiVs blas64.GeneralCols, psiVs blas64.GeneralCols, err error) C.SEXP {
	r := C.allocList(6)
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, 6)
	C.Rf_protect(names)
	arg := r
	C.SET_STRING_ELT(names, 0, C.Rf_mkCharLenCE(C._GoStringPtr("ccors"), 5, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Slice___float64(ccors))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 1, C.Rf_mkCharLenCE(C._GoStringPtr("pVecs"), 5, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Named_gonum_org_v1_gonum_blas_blas64_GeneralCols(pVecs))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 2, C.Rf_mkCharLenCE(C._GoStringPtr("qVecs"), 5, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Named_gonum_org_v1_gonum_blas_blas64_GeneralCols(qVecs))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 3, C.Rf_mkCharLenCE(C._GoStringPtr("phiVs"), 5, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Named_gonum_org_v1_gonum_blas_blas64_GeneralCols(phiVs))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 4, C.Rf_mkCharLenCE(C._GoStringPtr("psiVs"), 5, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Named_gonum_org_v1_gonum_blas_blas64_GeneralCols(psiVs))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 5, C.Rf_mkCharLenCE(C._GoStringPtr("err"), 3, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Named_error(err))
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	C.Rf_unprotect(2)
	return r
}

func unpackSEXP_types_Basic_int(p C.SEXP) int {
	return int(*C.INTEGER(p))
}

func unpackSEXP_types_Named_gonum_org_v1_gonum_blas_blas64_GeneralCols(p C.SEXP) blas64.GeneralCols {
	return unpackSEXP_types_Struct_struct_Rows_int__Cols_int__Data___float64__Stride_int_(p)
}

func unpackSEXP_types_Slice___float64(p C.SEXP) []float64 {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := C.Rf_xlength(p)
	return (*[70368744177664]float64)(unsafe.Pointer(C.REAL(p)))[:n:n]
}

func unpackSEXP_types_Struct_struct_Rows_int__Cols_int__Data___float64__Stride_int_(p C.SEXP) struct{Rows int; Cols int; Data []float64; Stride int} {
	switch n := C.Rf_xlength(p); {
	case n < 4:
		panic(`missing list element for struct{Rows int; Cols int; Data []float64; Stride int}`)
	case n > 4:
		err := C.CString(`extra list element ignored for struct{Rows int; Cols int; Data []float64; Stride int}`)
		C.R_error(err)
		C.free(unsafe.Pointer(err))
	}
	var r struct{Rows int; Cols int; Data []float64; Stride int}
	var i C.int
	key_Rows := C.CString("Rows")
	defer C.free(unsafe.Pointer(key_Rows))
	i = C.getListElementIndex(p, key_Rows)
	if i < 0 {
		panic("no list element for field: Rows")
	}
	r.Rows = unpackSEXP_types_Basic_int(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	key_Cols := C.CString("Cols")
	defer C.free(unsafe.Pointer(key_Cols))
	i = C.getListElementIndex(p, key_Cols)
	if i < 0 {
		panic("no list element for field: Cols")
	}
	r.Cols = unpackSEXP_types_Basic_int(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	key_Data := C.CString("Data")
	defer C.free(unsafe.Pointer(key_Data))
	i = C.getListElementIndex(p, key_Data)
	if i < 0 {
		panic("no list element for field: Data")
	}
	r.Data = unpackSEXP_types_Slice___float64(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	key_Stride := C.CString("Stride")
	defer C.free(unsafe.Pointer(key_Stride))
	i = C.getListElementIndex(p, key_Stride)
	if i < 0 {
		panic("no list element for field: Stride")
	}
	r.Stride = unpackSEXP_types_Basic_int(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	return r
}

func packSEXP_types_Basic_float64(p float64) C.SEXP {
	return C.ScalarReal(C.double(p))
}

func packSEXP_types_Basic_int(p int) C.SEXP {
	return C.ScalarInteger(C.int(p))
}

func packSEXP_types_Basic_string(p string) C.SEXP {
	s := C.Rf_mkCharLenCE(C._GoStringPtr(p), C.int(len(p)), C.CE_UTF8)
	return C.ScalarString(s)
}

func packSEXP_types_Named_error(p error) C.SEXP {
	if p == nil {
		return C.R_NilValue
	}
	return packSEXP_types_Basic_string(p.Error())
}

func packSEXP_types_Named_gonum_org_v1_gonum_blas_blas64_GeneralCols(p blas64.GeneralCols) C.SEXP {
	return packSEXP_types_Struct_struct_Rows_int__Cols_int__Data___float64__Stride_int_(p)
}

func packSEXP_types_Slice___float64(p []float64) C.SEXP {
	r := C.Rf_allocVector(C.REALSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	s := (*[70368744177664]float64)(unsafe.Pointer(C.REAL(r)))[:len(p):len(p)]
	copy(s, p)
	C.Rf_unprotect(1)
	return r
}

func packSEXP_types_Struct_struct_Rows_int__Cols_int__Data___float64__Stride_int_(p struct{Rows int; Cols int; Data []float64; Stride int}) C.SEXP {
	r := C.allocList(4)
	C.Rf_protect(r)
	names := C.Rf_allocVector(C.STRSXP, 4)
	C.Rf_protect(names)
	arg := r
	C.SET_STRING_ELT(names, 0, C.Rf_mkCharLenCE(C._GoStringPtr(`Rows`), 4, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_int(p.Rows))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 1, C.Rf_mkCharLenCE(C._GoStringPtr(`Cols`), 4, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_int(p.Cols))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 2, C.Rf_mkCharLenCE(C._GoStringPtr(`Data`), 4, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Slice___float64(p.Data))
	arg = C.CDR(arg)
	C.SET_STRING_ELT(names, 3, C.Rf_mkCharLenCE(C._GoStringPtr(`Stride`), 6, C.CE_UTF8))
	C.SETCAR(arg, packSEXP_types_Basic_int(p.Stride))
	C.setAttrib(r, packSEXP_types_Basic_string(`names`), names)
	C.Rf_unprotect(2)
	return r
}

func main() {}
