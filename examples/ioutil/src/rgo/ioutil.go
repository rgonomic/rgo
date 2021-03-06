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

	"io/ioutil"
)

//export Wrapped_ReadFile
func Wrapped_ReadFile(_R_filename C.SEXP) C.SEXP {
	defer func() {
		r := recover()
		if r != nil {
			err := C.CString(fmt.Sprint(r))
			C.R_error(err)
			C.free(unsafe.Pointer(err))
		}
	}()

	_p0 := unpackSEXP_types_Basic_string(_R_filename)
	_r0, _r1 := ioutil.ReadFile(_p0)
	return packSEXP_ReadFile(_r0, _r1)
}

func packSEXP_ReadFile(p0 []uint8, p1 error) C.SEXP {
	r := C.allocVector(C.VECSXP, 2)
	C.Rf_protect(r)
	defer C.Rf_unprotect(1)
	names := C.Rf_allocVector(C.STRSXP, 2)
	C.Rf_protect(names)
	defer C.Rf_unprotect(1)
	C.SET_STRING_ELT(names, 0, C.Rf_mkCharLenCE(C._GoStringPtr("r0"), 2, C.CE_UTF8))
	C.SET_VECTOR_ELT(r, 0, packSEXP_types_Slice___uint8(p0))
	C.SET_STRING_ELT(names, 1, C.Rf_mkCharLenCE(C._GoStringPtr("r1"), 2, C.CE_UTF8))
	C.SET_VECTOR_ELT(r, 1, packSEXP_types_Named_error(p1))
	C.setAttrib(r, C.R_NamesSymbol, names)
	return r
}

func unpackSEXP_types_Basic_string(p C.SEXP) string {
	return C.R_gostring(p, 0)
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

func packSEXP_types_Slice___uint8(p []uint8) C.SEXP {
	if p == nil {
		return C.R_NilValue
	}
	r := C.Rf_allocVector(C.RAWSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	defer C.Rf_unprotect(1)
	s := (*[562949953421312]uint8)(unsafe.Pointer(C.RAW(r)))[:len(p)]
	copy(s, p)
	return r
}

func main() {}
