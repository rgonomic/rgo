package sexp

import (
	"fmt"
	"unsafe"
)

// SEXP is the R SEXP type.
type SEXP struct {
	Header
	Data
}

// Header is the SEXPREC_HEADER #define in Rinternals.h.
type Header struct {
	Info   uint64 // Info is the sxpinfo_struct bitfield.
	Attrib unsafe.Pointer
	Next   unsafe.Pointer
	Prev   unsafe.Pointer
}

const (
	typeBits  = 5
	gpBits    = 16
	gcclsBits = 3
	namedBits = 16
	extraBits = 32 - namedBits
)

const (
	typ    = 0
	scalar = typ + typeBits // a numeric vector of length 1?
	obj    = scalar + 1     // an object with a class attribute?
	alt    = obj + 1        // an ALTREP object?
	gp     = alt + 1        // general purpose
	mark   = gp + 16        // mark object as ‘in use’ in GC
	debug  = mark + 1
	trace  = debug + 1 // functions and memory tracing
	spare  = trace + 1 // used on closures and when REFCNT is defined
	gcgen  = spare + 1 // old generation number
	gccls  = gcgen + 1 // node class
	named  = gccls + 3 // used to control copying
	extra  = named + namedBits
)

func mask(offset, bits int) uint64 {
	return ^uint64(0) << (64 - bits) >> (64 - bits - offset)
}

func (h Header) GoString() string {
	return fmt.Sprintf("R.SEXP{Type:%s, IsScalar:%t, IsObject:%t, IsAltRep:%t, GP:0b%016b, IsInUse:%t, Debug:%t, Trace:%t, Spare:%t, GCGeneration:%d, GCClass:%d, Named:%d, Extra:0b%016b}",
		SEXPType(h.Info&mask(typ, typeBits)),
		h.Info&mask(scalar, 1) != 0,
		h.Info&mask(obj, 1) != 0,
		h.Info&mask(alt, 1) != 0,
		(h.Info&mask(gp, gpBits))>>gp,
		h.Info&mask(mark, 1) != 0,
		h.Info&mask(debug, 1) != 0,
		h.Info&mask(trace, 1) != 0,
		h.Info&mask(spare, 1) != 0,
		h.Info&mask(gcgen, 1),
		(h.Info&mask(gccls, gcclsBits))>>gccls,
		(h.Info&mask(named, namedBits))>>named,
		(h.Info&mask(extra, extraBits))>>extra,
	)
}

// SEXPType is the SEXPTYPE enum in Rinternals.h
type SEXPType byte

//go:generate stringer -type=SEXPType
const (
	NILSXP     SEXPType = 0  // nil = NULL
	SYMSXP     SEXPType = 1  // symbols
	LISTSXP    SEXPType = 2  // lists of dotted pairs
	CLOSXP     SEXPType = 3  // closures
	ENVSXP     SEXPType = 4  // environments
	PROMSXP    SEXPType = 5  // promises: [un]evaluated closure arguments
	LANGSXP    SEXPType = 6  // language constructs (special lists)
	SPECIALSXP SEXPType = 7  // special forms
	BUILTINSXP SEXPType = 8  // builtin non-special forms
	CHARSXP    SEXPType = 9  // "scalar" string type (internal only)
	LGLSXP     SEXPType = 10 // logical vectors
	INTSXP     SEXPType = 13 // integer vectors
	REALSXP    SEXPType = 14 // real variables
	CPLXSXP    SEXPType = 15 // complex variables
	STRSXP     SEXPType = 16 // string vectors
	DOTSXP     SEXPType = 17 // dot-dot-dot object
	ANYSXP     SEXPType = 18 // make "any" args work
	VECSXP     SEXPType = 19 // generic vectors
	EXPRSXP    SEXPType = 20 // expressions vectors
	BCODESXP   SEXPType = 21 // byte code
	EXTPTRSXP  SEXPType = 22 // external pointer
	WEAKREFSXP SEXPType = 23 // weak reference
	RAWSXP     SEXPType = 24 // raw bytes
	S4SXP      SEXPType = 25 // S4 non-vector

	NEWSXP  SEXPType = 30 // fresh node creaed in new page
	FREESXP SEXPType = 31 // node released by GC

	FUNSXP SEXPType = 99 // Closure or Builtin

)

type Data struct {
	_ unsafe.Pointer
	_ unsafe.Pointer
	_ unsafe.Pointer
}

// PrintSEXPHeader prints the header of the SEXP value passed to it in
// Go syntax and returns it unaltered.
func PrintSEXPHeader(sexp unsafe.Pointer) unsafe.Pointer {
	p := (*SEXP)(sexp)
	fmt.Printf("%#v\n", p.Header)
	return sexp
}
