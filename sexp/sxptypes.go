// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build generate

package sexp

/*
#cgo pkg-config: libR

#define USE_RINTERNALS
#include <R.h>
#include <Rinternals.h>
*/
import "C"

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
