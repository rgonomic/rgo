// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sexp

import "unsafe"

/*
#cgo pkg-config: libR

#define USE_RINTERNALS
#include <R.h>
#include <Rinternals.h>

typedef struct sxpinfo_struct sxpinfo_struct;

typedef struct SEXPREC sexprec;

typedef struct VECTOR_SEXPREC vector_sexprec;

typedef struct prim_sexprec {
    SEXPREC_HEADER;
    struct primsxp_struct prim_sxp;
} prim_sexprec;

typedef struct sym_sexprec {
    SEXPREC_HEADER;
    struct symsxp_struct sym_sxp;
} sym_sexprec;

typedef struct list_sexprec {
    SEXPREC_HEADER;
    struct listsxp_struct list_sxp;
} list_sexprec;

typedef struct env_sexprec {
    SEXPREC_HEADER;
    struct envsxp_struct env_sxp;
} env_sexprec;

typedef struct clo_sexprec {
    SEXPREC_HEADER;
    struct closxp_struct clos_sxp;
} clo_sexprec;

typedef struct prom_sexprec {
    SEXPREC_HEADER;
    struct promsxp_struct prom_sxp;
} prom_sexprec;
*/
import "C"

type sxpinfo C.sxpinfo_struct

type sexprec C.sexprec

type vector_sexprec C.vector_sexprec

type prim_sexprec C.prim_sexprec

type sym_sexprec C.sym_sexprec

type list_sexprec C.list_sexprec

type env_sexprec C.env_sexprec

type clo_sexprec C.clo_sexprec

type prom_sexprec C.prom_sexprec

const (
	NILSXP     Type = C.NILSXP     // nil = NULL
	SYMSXP     Type = C.SYMSXP     // symbols
	LISTSXP    Type = C.LISTSXP    // lists of dotted pairs
	CLOSXP     Type = C.CLOSXP     // closures
	ENVSXP     Type = C.ENVSXP     // environments
	PROMSXP    Type = C.PROMSXP    // promises: [un]evaluated closure arguments
	LANGSXP    Type = C.LANGSXP    // language constructs (special lists)
	SPECIALSXP Type = C.SPECIALSXP // special forms
	BUILTINSXP Type = C.BUILTINSXP // builtin non-special forms
	CHARSXP    Type = C.CHARSXP    // "scalar" string type (internal only)
	LGLSXP     Type = C.LGLSXP     // logical vectors
	INTSXP     Type = C.INTSXP     // integer vectors
	REALSXP    Type = C.REALSXP    // real variables
	CPLXSXP    Type = C.CPLXSXP    // complex variables
	STRSXP     Type = C.STRSXP     // string vectors
	DOTSXP     Type = C.DOTSXP     // dot-dot-dot object
	ANYSXP     Type = C.ANYSXP     // make "any" args work
	VECSXP     Type = C.VECSXP     // generic vectors
	EXPRSXP    Type = C.EXPRSXP    // expressions vectors
	BCODESXP   Type = C.BCODESXP   // byte code
	EXTPTRSXP  Type = C.EXTPTRSXP  // external pointer
	WEAKREFSXP Type = C.WEAKREFSXP // weak reference
	RAWSXP     Type = C.RAWSXP     // raw bytes
	S4SXP      Type = C.S4SXP      // S4 non-vector
	NEWSXP     Type = C.NEWSXP     // fresh node created in new page
	FREESXP    Type = C.FREESXP    // node released by GC
	FUNSXP     Type = C.FUNSXP     // Closure or Builtin
)

var (
	// NilValue is the R_NilValue as a *Value type. It must not be altered.
	NilValue *Value = asValue(C.R_NilValue)

	// UnboundValue is the R_UnboundValue as a *Value type. It must not be altered.
	UnboundValue *Value = asValue(C.R_UnboundValue)
)

var r_NaInt = int32(C.R_NaInt)

func asValue(sexp C.SEXP) *Value {
	return (*Value)(unsafe.Pointer(sexp))
}

func allocateList(n int) unsafe.Pointer {
	return unsafe.Pointer(C.Rf_allocList(C.int(n)))
}

func allocateString(s string) unsafe.Pointer {
	return unsafe.Pointer(C.Rf_mkCharLenCE(C._GoStringPtr(s), C.int(len(s)), C.CE_UTF8))
}

func allocateStringFromBytes(b []byte) unsafe.Pointer {
	// This makes use of the pointer being the first field of the slice header.
	p := *(**C.char)(unsafe.Pointer(&b))
	return unsafe.Pointer(C.Rf_mkCharLenCE(p, C.int(len(b)), C.CE_UTF8))
}

func allocateVector(typ Type, n int) unsafe.Pointer {
	return unsafe.Pointer(C.Rf_allocVector(C.SEXPTYPE(typ), C.R_xlen_t(n)))
}

func protect(sexp unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.Rf_protect((C.SEXP)(sexp)))
}

func unprotect(n int) {
	C.Rf_unprotect(C.int(n))
}
