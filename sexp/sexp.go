// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sexp

import (
	"fmt"
	"unsafe"
)

// Info corresponds to the sxpinfo_struct type defined in Rinternals.h.
type Info uint64

// An "invalid array index" compiler error after go generate means that
// the R sxpinfo_struct type definition has changed.
// If this happens, definition of Info needs to be brought back in line
// with sxpinfo_struct in Rinternals.h.
var _ = [1]struct{}{}[int(unsafe.Sizeof(sxpinfo{}))-int(unsafe.Sizeof(Info(0)))]

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

func (i Info) String() string {
	return fmt.Sprintf("Info{Type:%s, IsScalar:%t, IsObject:%t, IsAltRep:%t, GP:0b%016b, IsInUse:%t, Debug:%t, Trace:%t, Spare:%t, GCGeneration:%d, GCClass:%d, Named:%d, Extra:0b%016b}",
		i.Type(),
		i.IsScalar(),
		i.IsObject(),
		i.IsAltRep(),
		i.GP(),
		i.IsInUse(),
		i.Debug(),
		i.Trace(),
		i.Space(),
		i.GCGen(),
		i.GCClass(),
		i.Named(),
		i.Extra(),
	)
}

func (i Info) Type() SEXPType { return SEXPType(i.mask(typ, typeBits)) }
func (i Info) IsScalar() bool { return i.mask(scalar, 1) != 0 }
func (i Info) IsObject() bool { return i.mask(obj, 1) != 0 }
func (i Info) IsAltRep() bool { return i.mask(alt, 1) != 0 }
func (i Info) GP() uint16     { return uint16(i.mask(gp, gpBits)) }
func (i Info) IsInUse() bool  { return i.mask(mark, 1) != 0 }
func (i Info) Debug() bool    { return i.mask(debug, 1) != 0 }
func (i Info) Trace() bool    { return i.mask(trace, 1) != 0 }
func (i Info) Space() bool    { return i.mask(spare, 1) != 0 }
func (i Info) GCGen() int     { return int(i.mask(gcgen, 1)) }
func (i Info) GCClass() int   { return int(i.mask(gccls, gcclsBits)) }
func (i Info) Named() uint16  { return uint16(i.mask(named, namedBits)) }
func (i Info) Extra() uint16  { return uint16(i.mask(extra, extraBits)) }

func (i Info) mask(offset, bits int) uint64 {
	return (uint64(i) >> offset) & (^uint64(0) >> (64 - bits))
}

// SEXPType is the SEXPTYPE enum defined in Rinternals.h.
//go:generate stringer -type=SEXPType
type SEXPType byte

// Value is an SEXP value.
type Value struct {
	sexprec
}

// Info returns the information field of the SEXP value.
func (v *Value) Info() Info {
	return *(*Info)(unsafe.Pointer(&v.Sxpinfo))
}

// Attributes returns the attributes of the SEXP value.
func (v *sexprec) Attributes() *Value {
	return (*Value)(unsafe.Pointer(&v.Attrib))
}

// Interface returns a Go value corresponding to the SEXP type specified
// in the SEXP info field.
func (v *Value) Interface() interface{} {
	switch typ := v.Info().Type(); typ {
	case NILSXP:
		return v
	case SYMSXP:
		return (*sym_sexprec)(unsafe.Pointer(v))
	case LISTSXP:
		return (*list_sexprec)(unsafe.Pointer(v))
	case CLOSXP:
		return (*clo_sexprec)(unsafe.Pointer(v))
	case ENVSXP:
		return (*env_sexprec)(unsafe.Pointer(v))
	case PROMSXP:
		return (*prom_sexprec)(unsafe.Pointer(v))
	case LANGSXP:
		return (*list_sexprec)(unsafe.Pointer(v))
	case SPECIALSXP:
		return (*prim_sexprec)(unsafe.Pointer(v))
	case BUILTINSXP:
		return (*prim_sexprec)(unsafe.Pointer(v))
	case CHARSXP:
		return (*Character)(unsafe.Pointer(v))
	case LGLSXP:
		return (*Logical)(unsafe.Pointer(v))
	case INTSXP:
		return (*Integer)(unsafe.Pointer(v))
	case REALSXP:
		return (*Real)(unsafe.Pointer(v))
	case CPLXSXP:
		return (*Complex)(unsafe.Pointer(v))
	case STRSXP:
		return (*String)(unsafe.Pointer(v))
	case DOTSXP:
		return (*list_sexprec)(unsafe.Pointer(v))
	case ANYSXP:
		return v
	case VECSXP:
		return (*Vector)(unsafe.Pointer(v))
	case EXPRSXP:
		return (*vector_sexprec)(unsafe.Pointer(v))
	case BCODESXP:
		return v
	case EXTPTRSXP:
		return v
	case WEAKREFSXP:
		return (*vector_sexprec)(unsafe.Pointer(v))
	case RAWSXP:
		return (*Raw)(unsafe.Pointer(v))
	case S4SXP:
		return v
	case NEWSXP:
		return v
	case FREESXP:
		return v
	case FUNSXP:
		return v
	default:
		panic(fmt.Sprintf("unhandled SEXPTYPE: 0x%x", typ))
	}
}

// Len returns the number of elements in the vector.
func (v *vector_sexprec) Len() int {
	return int(v.Vecsxp.Length)
}

// base returns the address of the first element of the vector.
func (v *vector_sexprec) base() unsafe.Pointer {
	return add(unsafe.Pointer(v), unsafe.Sizeof(vector_sexprec{}))
}

func add(addr unsafe.Pointer, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(addr) + offset)
}

// Attributes returns the attributes of the SEXP value.
func (v *vector_sexprec) Attributes() *Value {
	return (*Value)(unsafe.Pointer(&v.Attrib))
}

// Integer is an R integer vector.
type Integer struct {
	vector_sexprec
}

// Vector returns a slice corresponding to the R vector.
func (v *Integer) Vector() []int32 {
	n := v.Len()
	return (*[1 << 30]int32)(v.base())[:n:n]
}

// Logical is an R logical vector.
type Logical struct {
	vector_sexprec
}

// Vector returns a slice corresponding to the R vector.
func (v *Logical) Vector() []int32 {
	n := v.Len()
	return (*[1 << 30]int32)(v.base())[:n:n]
}

// Real is an R real vector.
type Real struct {
	vector_sexprec
}

// Vector returns a slice corresponding to the R vector.
func (v *Real) Vector() []float64 {
	n := v.Len()
	return (*[1 << 30]float64)(v.base())[:n:n]
}

// Complex is an R complex vector.
type Complex struct {
	vector_sexprec
}

// Vector returns a slice corresponding to the R vector.
func (v *Complex) Vector() []complex128 {
	n := v.Len()
	return (*[1 << 30]complex128)(v.base())[:n:n]
}

// String is an R character vector.
type String struct {
	vector_sexprec
}

// Vector returns a slice corresponding to the R vector.
func (v *String) Vector() []*Character {
	n := v.Len()
	return (*[1 << 30]*Character)(v.base())[:n:n]
}

// Character is the R representation of a string.
type Character struct {
	vector_sexprec
}

// Bytes returns the bytes held by the R SEXP value.
func (v *Character) Bytes() []byte {
	n := v.Len()
	return (*[1 << 30]byte)(v.base())[:n:n]
}

// String returns a Go string corresponding the the R characters.
// The returned string is allocated by the Go runtime.
func (v *Character) String() string {
	return string(v.Bytes())
}

// Raw is an R raw vector.
type Raw struct {
	vector_sexprec
}

// Bytes returns the bytes held by the R SEXP value.
func (v *Raw) Bytes() []byte {
	n := v.Len()
	return (*[1 << 30]byte)(v.base())[:n:n]
}

// Vector is a generic R vector.
type Vector struct {
	vector_sexprec
}

// Vector returns a slice corresponding to the R vector.
func (v *Vector) Vector() []*Value {
	n := v.Len()
	return (*[1 << 30]*Value)(v.base())[:n:n]
}
