//go:generate bash -c "go tool cgo -godefs -- $(pkg-config --cflags libR) generate.go | gofmt > cgo_types.go"
//go:generate rm -rf _obj

// +build generate

package sexp

/*
#define USE_RINTERNALS
#include <R.h>
#include <Rinternals.h>

typedef struct sxpinfo_struct sxpinfo_struct;

typedef struct SEXPREC sexprec;

typedef struct VECTOR_SEXPREC vector_sexprec;

typedef struct vecsxp_struct  vecsxp_struct;

typedef struct primsxp_struct primsxp_struct;
typedef struct prim_sexprec {
    SEXPREC_HEADER;
    struct primsxp_struct prim_sxp;
} prim_sexprec;

typedef struct symsxp_struct  symsxp_struct;
typedef struct sym_sexprec {
    SEXPREC_HEADER;
    struct symsxp_struct sym_sxp;
} sym_sexprec;

typedef struct listsxp_struct listsxp_struct;
typedef struct list_sexprec {
    SEXPREC_HEADER;
    struct listsxp_struct list_sxp;
} list_sexprec;

typedef struct envsxp_struct  envsxp_struct;
typedef struct env_sexprec {
    SEXPREC_HEADER;
    struct envsxp_struct env_sxp;
} env_sexprec;

typedef struct closxp_struct  closxp_struct;
typedef struct clo_sexprec {
    SEXPREC_HEADER;
    struct closxp_struct clos_sxp;
} clo_sexprec;

typedef struct promsxp_struct promsxp_struct;
typedef struct prom_sexprec {
    SEXPREC_HEADER;
    struct promsxp_struct prom_sxp;
} prom_sexprec;
*/
import "C"

type sxpinfo C.sxpinfo_struct

type sexprec C.sexprec

type (
	vector_sexprec C.vector_sexprec
	vecsxp         C.vecsxp_struct
)

type (
	prim_sexprec C.prim_sexprec
	primsxp      C.primsxp_struct
)

type (
	sym_sexprec C.sym_sexprec
	symsxp      C.symsxp_struct
)

type (
	list_sexprec C.list_sexprec
	listsxp      C.listsxp_struct
)

type (
	env_sexprec C.env_sexprec
	envsxp      C.envsxp_struct
)

type (
	clo_sexprec C.clo_sexprec
	closxp      C.closxp_struct
)

type (
	prom_sexprec C.prom_sexprec
	promsxp      C.promsxp_struct
)

const (
	NILSXP     SEXPType = C.NILSXP     // nil = NULL
	SYMSXP     SEXPType = C.SYMSXP     // symbols
	LISTSXP    SEXPType = C.LISTSXP    // lists of dotted pairs
	CLOSXP     SEXPType = C.CLOSXP     // closures
	ENVSXP     SEXPType = C.ENVSXP     // environments
	PROMSXP    SEXPType = C.PROMSXP    // promises: [un]evaluated closure arguments
	LANGSXP    SEXPType = C.LANGSXP    // language constructs (special lists)
	SPECIALSXP SEXPType = C.SPECIALSXP // special forms
	BUILTINSXP SEXPType = C.BUILTINSXP // builtin non-special forms
	CHARSXP    SEXPType = C.CHARSXP    // "scalar" string type (internal only
	LGLSXP     SEXPType = C.LGLSXP     // logical vectors
	INTSXP     SEXPType = C.INTSXP     // integer vectors
	REALSXP    SEXPType = C.REALSXP    // real variables
	CPLXSXP    SEXPType = C.CPLXSXP    // complex variables
	STRSXP     SEXPType = C.STRSXP     // string vectors
	DOTSXP     SEXPType = C.DOTSXP     // dot-dot-dot object
	ANYSXP     SEXPType = C.ANYSXP     // make "any" args work
	VECSXP     SEXPType = C.VECSXP     // generic vectors
	EXPRSXP    SEXPType = C.EXPRSXP    // expressions vectors
	BCODESXP   SEXPType = C.BCODESXP   // byte code
	EXTPTRSXP  SEXPType = C.EXTPTRSXP  // external pointer
	WEAKREFSXP SEXPType = C.WEAKREFSXP // weak reference
	RAWSXP     SEXPType = C.RAWSXP     // raw bytes
	S4SXP      SEXPType = C.S4SXP      // S4 non-vector
	NEWSXP     SEXPType = C.NEWSXP     // fresh node created in new page
	FREESXP    SEXPType = C.FREESXP    // node released by GC
	FUNSXP     SEXPType = C.FUNSXP     // Closure or Builtin
)
