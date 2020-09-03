// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sexp

import (
	"bytes"
	"fmt"
	"unsafe"
)

// IsNull returned whether the Value is the R NULL value.
func IsNull(v *Value) bool {
	return v.Info().Type() == NILSXP
}

// Names returns the names of the Value.
func Names(v *Value) *String {
	if v, ok := v.Interface().(*List); ok {
		tags := v.tags()
		c := NewString(len(tags)).Protect()
		defer c.Unprotect()
		vec := c.Vector()
		for i := range vec {
			vec[i] = NewCharacter(tags[i])
		}
		return c
	}

	attr := v.Attributes()
	if IsNull(attr) {
		return nil
	}
	names := attr.Interface().(*List).Get([]byte("names"))
	found, ok := names.Interface().(*List)
	if !ok {
		return nil
	}
	return found.Head().Interface().(*String)
}

// Info corresponds to the sxpinfo_struct type defined in Rinternals.h.
type Info uint64

// An "invalid array index" compiler error means that the R sxpinfo_struct
// type definition has changed.
// If this happens, definition of Info needs to be brought back in line
// with sxpinfo_struct in Rinternals.h.
//
// This changed from 32 bits in R commit 14db4328 (svn revision 73243),
// so we do not work with R code prior to that.
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

func (i Info) Type() Type     { return Type(i.mask(typ, typeBits)) }
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

//go:generate stringer -type=Type

// Type is the SEXPTYPE enum defined in Rinternals.h.
type Type byte

// Info returns the information field of the SEXP value.
func (v *sexprec) Info() Info {
	return *(*Info)(unsafe.Pointer(&v.sxpinfo))
}

// Value returns the generic state of the SEXP value.
func (v *sexprec) Value() *Value {
	return (*Value)(unsafe.Pointer(v))
}

// Attributes returns the attributes of the SEXP value.
func (v *sexprec) Attributes() *Value {
	return (*Value)(unsafe.Pointer(v.attrib))
}

// Pointer returns an unsafe pointer to the SEXP value.
func (v *sexprec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

// Info returns the information field of the SEXP value.
func (v *vector_sexprec) Info() Info {
	return *(*Info)(unsafe.Pointer(&v.sxpinfo))
}

// Value returns the generic state of the SEXP value.
func (v *vector_sexprec) Value() *Value {
	return (*Value)(unsafe.Pointer(v))
}

// Attributes returns the attributes of the SEXP value.
func (v *vector_sexprec) Attributes() *Value {
	return (*Value)(unsafe.Pointer(v.attrib))
}

// Pointer returns an unsafe pointer to the SEXP value.
func (v *vector_sexprec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

// Info returns the information field of the SEXP value.
func (v *list_sexprec) Info() Info {
	return *(*Info)(unsafe.Pointer(&v.sxpinfo))
}

// Value returns the generic state of the SEXP value.
func (v *list_sexprec) Value() *Value {
	return (*Value)(unsafe.Pointer(v))
}

// Attributes returns the attributes of the SEXP value.
func (v *list_sexprec) Attributes() *Value {
	return (*Value)(unsafe.Pointer(v.attrib))
}

// Pointer returns an unsafe pointer to the SEXP value.
func (v *list_sexprec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

// Info returns the information field of the SEXP value.
func (v *env_sexprec) Info() Info {
	return *(*Info)(unsafe.Pointer(&v.sxpinfo))
}

// Value returns the generic state of the SEXP value.
func (v *env_sexprec) Value() *Value {
	return (*Value)(unsafe.Pointer(v))
}

// Attributes returns the attributes of the SEXP value.
func (v *env_sexprec) Attributes() *Value {
	return (*Value)(unsafe.Pointer(v.attrib))
}

// Pointer returns an unsafe pointer to the SEXP value.
func (v *env_sexprec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

// Info returns the information field of the SEXP value.
func (v *prom_sexprec) Info() Info {
	return *(*Info)(unsafe.Pointer(&v.sxpinfo))
}

// Value returns the generic state of the SEXP value.
func (v *prom_sexprec) Value() *Value {
	return (*Value)(unsafe.Pointer(v))
}

// Attributes returns the attributes of the SEXP value.
func (v *prom_sexprec) Attributes() *Value {
	return (*Value)(unsafe.Pointer(v.attrib))
}

// Pointer returns an unsafe pointer to the SEXP value.
func (v *prom_sexprec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

// Info returns the information field of the SEXP value.
func (v *clo_sexprec) Info() Info {
	return *(*Info)(unsafe.Pointer(&v.sxpinfo))
}

// Value returns the generic state of the SEXP value.
func (v *clo_sexprec) Value() *Value {
	return (*Value)(unsafe.Pointer(v))
}

// Attributes returns the attributes of the SEXP value.
func (v *clo_sexprec) Attributes() *Value {
	return (*Value)(unsafe.Pointer(v.attrib))
}

// Info returns the information field of the SEXP value.
func (v *prim_sexprec) Info() Info {
	return *(*Info)(unsafe.Pointer(&v.sxpinfo))
}

// Value returns the generic state of the SEXP value.
func (v *prim_sexprec) Value() *Value {
	return (*Value)(unsafe.Pointer(v))
}

// Attributes returns the attributes of the SEXP value.
func (v *prim_sexprec) Attributes() *Value {
	return (*Value)(unsafe.Pointer(v.attrib))
}

// Pointer returns an unsafe pointer to the SEXP value.
func (v *prim_sexprec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

// Info returns the information field of the SEXP value.
func (v *sym_sexprec) Info() Info {
	return *(*Info)(unsafe.Pointer(&v.sxpinfo))
}

// Value returns the generic state of the SEXP value.
func (v *sym_sexprec) Value() *Value {
	return (*Value)(unsafe.Pointer(v))
}

// Attributes returns the attributes of the SEXP value.
func (v *sym_sexprec) Attributes() *Value {
	return (*Value)(unsafe.Pointer(v.attrib))
}

// Pointer returns an unsafe pointer to the SEXP value.
func (v *sym_sexprec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

// Value is an SEXP value.
type Value struct {
	sexprec
}

// Interface returns a Go value corresponding to the SEXP type specified
// in the SEXP info field.
func (v *Value) Interface() interface{} {
	switch typ := v.Info().Type(); typ {
	case NILSXP:
		return v
	case SYMSXP:
		return (*Symbol)(unsafe.Pointer(v))
	case LISTSXP:
		return (*List)(unsafe.Pointer(v))
	case CLOSXP:
		return (*Closure)(unsafe.Pointer(v))
	case ENVSXP:
		return (*Environment)(unsafe.Pointer(v))
	case PROMSXP:
		return (*Promise)(unsafe.Pointer(v))
	case LANGSXP:
		return (*Lang)(unsafe.Pointer(v))
	case SPECIALSXP:
		return (*Special)(unsafe.Pointer(v))
	case BUILTINSXP:
		return (*Builtin)(unsafe.Pointer(v))
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
		return (*Dot)(unsafe.Pointer(v))
	case ANYSXP:
		return v
	case VECSXP:
		return (*Vector)(unsafe.Pointer(v))
	case EXPRSXP:
		return (*Expression)(unsafe.Pointer(v))
	case BCODESXP:
		return v
	case EXTPTRSXP:
		return v
	case WEAKREFSXP:
		return (*WeakReference)(unsafe.Pointer(v))
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
	return int(v.vecsxp.length)
}

// base returns the address of the first element of the vector.
func (v *vector_sexprec) base() unsafe.Pointer {
	return add(unsafe.Pointer(v), unsafe.Sizeof(vector_sexprec{}))
}

func add(addr unsafe.Pointer, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(addr) + offset)
}

// Integer is an R integer vector.
type Integer struct {
	vector_sexprec
}

// NewInteger returns an integer vector with length n.
//
// The allocation is made by the R runtime. The returned value may need to
// call its Protect method.
func NewInteger(n int) *Integer {
	return (*Integer)(allocateVector(INTSXP, n))
}

// Protect protects the SEXP value and returns it.
func (v *Integer) Protect() *Integer {
	return (*Integer)(protect(unsafe.Pointer(v)))
}

// Unprotect unprotects the SEXP. It is equivalent to UNPROTECT(1).
func (v *Integer) Unprotect() {
	unprotect(1)
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

// NewLogical returns a logical vector with length n.
//
// The allocation is made by the R runtime. The returned value may need to
// call its Protect method.
func NewLogical(n int) *Logical {
	return (*Logical)(allocateVector(LGLSXP, n))
}

// Protect protects the SEXP value and returns it.
func (v *Logical) Protect() *Logical {
	return (*Logical)(protect(unsafe.Pointer(v)))
}

// Unprotect unprotects the SEXP. It is equivalent to UNPROTECT(1).
func (v *Logical) Unprotect() {
	unprotect(1)
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

// NewReal returns a real vector with length n.
//
// The allocation is made by the R runtime. The returned value may need to
// call its Protect method.
func NewReal(n int) *Real {
	return (*Real)(allocateVector(REALSXP, n))
}

// Protect protects the SEXP value and returns it.
func (v *Real) Protect() *Real {
	return (*Real)(protect(unsafe.Pointer(v)))
}

// Unprotect unprotects the SEXP. It is equivalent to UNPROTECT(1).
func (v *Real) Unprotect() {
	unprotect(1)
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

// NewComplex returns a complex vector with length n.
//
// The allocation is made by the R runtime. The returned value may need to
// call its Protect method.
func NewComplex(n int) *Complex {
	return (*Complex)(allocateVector(CPLXSXP, n))
}

// Protect protects the SEXP value and returns it.
func (v *Complex) Protect() *Complex {
	return (*Complex)(protect(unsafe.Pointer(v)))
}

// Unprotect unprotects the SEXP. It is equivalent to UNPROTECT(1).
func (v *Complex) Unprotect() {
	unprotect(1)
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

// NewString returns a character vector with length n.
//
// The allocation is made by the R runtime. The returned value may need to
// call its Protect method.
func NewString(n int) *String {
	return (*String)(allocateVector(STRSXP, n))
}

// Protect protects the SEXP value and returns it.
func (v *String) Protect() *String {
	return (*String)(protect(unsafe.Pointer(v)))
}

// Unprotect unprotects the SEXP. It is equivalent to UNPROTECT(1).
func (v *String) Unprotect() {
	unprotect(1)
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

// NewCharacter returns a scalar string corresponding to s.
//
// The allocation is made by the R runtime. The returned value may need to
// call its Protect method.
func NewCharacter(s string) *Character {
	return (*Character)(allocateString(s))
}

// Protect protects the SEXP value and returns it.
func (v *Character) Protect() *Character {
	return (*Character)(protect(unsafe.Pointer(v)))
}

// Unprotect unprotects the SEXP. It is equivalent to UNPROTECT(1).
func (v *Character) Unprotect() {
	unprotect(1)
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

// NewRaw returns a raw vector with length n.
//
// The allocation is made by the R runtime. The returned value may need to
// call its Protect method.
func NewRaw(n int) *Raw {
	return (*Raw)(allocateVector(RAWSXP, n))
}

// Protect protects the SEXP value and returns it.
func (v *Raw) Protect() *Raw {
	return (*Raw)(protect(unsafe.Pointer(v)))
}

// Unprotect unprotects the SEXP. It is equivalent to UNPROTECT(1).
func (v *Raw) Unprotect() {
	unprotect(1)
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

// NewVector returns a generic vector with length n.
//
// The allocation is made by the R runtime. The returned value may need to
// call its Protect method.
func NewVector(n int) *Vector {
	return (*Vector)(allocateVector(VECSXP, n))
}

// Protect protects the SEXP value and returns it.
func (v *Vector) Protect() *Vector {
	return (*Vector)(protect(unsafe.Pointer(v)))
}

// Unprotect unprotects the SEXP. It is equivalent to UNPROTECT(1).
func (v *Vector) Unprotect() {
	unprotect(1)
}

// Vector returns a slice corresponding to the R vector.
func (v *Vector) Vector() []*Value {
	n := v.Len()
	return (*[1 << 30]*Value)(v.base())[:n:n]
}

// Expression is an R expression.
type Expression struct {
	vector_sexprec
}

// Expression returns a slice corresponding to the R expression.
func (v *Expression) Vector() []*Value {
	n := v.Len()
	return (*[1 << 30]*Value)(v.base())[:n:n]
}

// WeakReference is an R weak reference.
type WeakReference struct {
	vector_sexprec
}

// WeakReference returns the four elements of an R weak reference.
func (v *WeakReference) Vector() []*Value {
	n := v.Len()
	return (*[1 << 30]*Value)(v.base())[:n:n]
}

// List is an R linked list.
type List struct {
	list_sexprec
}

// NewList returns a list with length n.
func NewList(n int) *List {
	return (*List)(allocateList(n))
}

// Protect protects the SEXP value and returns it.
func (v *List) Protect() *List {
	return (*List)(protect(unsafe.Pointer(v)))
}

// Unprotect unprotects the SEXP. It is equivalent to UNPROTECT(1).
func (v *List) Unprotect() {
	unprotect(1)
}

// Head returns the first element of the list (CAR/lisp in R terminology).
func (v *List) Head() *Value {
	return (*Value)(unsafe.Pointer(v.list_sxp.carval))
}

// Tail returns the remaining elements of the list (CDR/lisp in R terminology).
func (v *List) Tail() *Value {
	return (*Value)(unsafe.Pointer(v.list_sxp.cdrval))
}

// Tag returns the list's tag value.
func (v *List) Tag() *Value {
	return (*Value)(unsafe.Pointer(v.list_sxp.tagval))
}

// Get returns the the Value associated with the given tag in the list.
func (v *List) Get(tag []byte) *Value {
	curr := v
	for !IsNull(curr.Value()) {
		t := curr.Tag().Value()
		if !IsNull(t) {
			sym := t.Interface().(*Symbol)
			if bytes.Equal(sym.Name().Bytes(), tag) {
				return (*Value)(curr.Pointer())
			}
		}
		tail := curr.Tail()
		if tail, ok := tail.Value().Interface().(*List); ok {
			curr = (*List)(tail.Pointer())
			continue
		}
		break
	}
	return NilValue
}

// tags returns all the tags for the list.
func (v *List) tags() []string {
	var tags []string
	curr := v
	for !IsNull(curr.Value()) {
		t := curr.Tag().Value()
		if !IsNull(t) {
			tag := t.Interface().(*Symbol).String()
			tags = append(tags, tag)
		}
		tail := curr.Tail()
		if tail, ok := tail.Value().Interface().(*List); ok {
			curr = (*List)(tail.Pointer())
			continue
		}
		break
	}
	return tags
}

// Lang is an R language object.
type Lang struct {
	list_sexprec
}

// Head returns the first element of the list (CAR/lisp in R terminology).
func (v *Lang) Head() *Value {
	return (*Value)(unsafe.Pointer(v.list_sxp.carval))
}

// Tail returns the remaining elements of the list (CDR/lisp in R terminology).
func (v *Lang) Tail() *Value {
	return (*Value)(unsafe.Pointer(v.list_sxp.cdrval))
}

// Tag returns the object's tag value.
func (v *Lang) Tag() *Value {
	return (*Value)(unsafe.Pointer(v.list_sxp.tagval))
}

// Dot is an R pairlist of promises.
type Dot struct {
	list_sexprec
}

// Head returns the first element of the list (CAR/lisp in R terminology).
func (v *Dot) Head() *Value {
	return (*Value)(unsafe.Pointer(v.list_sxp.carval))
}

// Tail returns the remaining elements of the list (CDR/lisp in R terminology).
func (v *Dot) Tail() *Value {
	return (*Value)(unsafe.Pointer(v.list_sxp.cdrval))
}

// Tag returns the object's tag value.
func (v *Dot) Tag() *Value {
	return (*Value)(unsafe.Pointer(v.list_sxp.tagval))
}

// Symbol is an R name value.
type Symbol struct {
	sym_sexprec
}

// Value returns the value of the symbol.
func (v *Symbol) SymbolValue() *Value {
	return (*Value)(unsafe.Pointer(v.sym_sxp.value))
}

// Name returns the name of the symbol
func (v *Symbol) Name() *Character {
	return (*Character)(unsafe.Pointer(v.sym_sxp.pname))
}

// String returns a Go string of the symbol name.
// The returned string is allocated by the Go runtime.
func (v *Symbol) String() string {
	return v.Name().String()
}

// Internal returns a pointer if the symbol is a .Internal function.
func (v *Symbol) Internal() *Value {
	return (*Value)(unsafe.Pointer(v.sym_sxp.internal))
}

// Promise is an R promise.
type Promise struct {
	prom_sexprec
}

// Value is value of the promise.
func (v *Promise) Value() *Value {
	return (*Value)(unsafe.Pointer(v.prom_sxp.value))
}

// Expression is the expression to be evaluated.
func (v *Promise) Expression() *Value {
	return (*Value)(unsafe.Pointer(v.prom_sxp.expr))
}

// Environment returns the environment in which to evaluate the expression.
func (v *Promise) Environment() *Value {
	return (*Value)(unsafe.Pointer(v.prom_sxp.env))
}

// Closure is an R closure.
type Closure struct {
	clo_sexprec
}

// Formals returns the formal arguments of the function.
func (v *Closure) Formals() *Value {
	return (*Value)(unsafe.Pointer(v.clos_sxp.formals))
}

// Body returns the body of the function.
func (v *Closure) Body() *Value {
	return (*Value)(unsafe.Pointer(v.clos_sxp.body))
}

// Environment returns the environment in which to evaluate the function.
func (v *Closure) Environment() *Value {
	return (*Value)(unsafe.Pointer(v.clos_sxp.env))
}

// Environment is a current execution environment.
type Environment struct {
	env_sexprec
}

// Frame returns the current frame.
func (v *Environment) Frame() *Value {
	return (*Value)(unsafe.Pointer(v.env_sxp.frame))
}

// Enclosing returns the enclosing environment.
func (v *Environment) Enclosing() *Value {
	return (*Value)(unsafe.Pointer(v.env_sxp.enclos))
}

// HashTable returns the environment's hash table.
func (v *Environment) HashTable() *Value {
	return (*Value)(unsafe.Pointer(v.env_sxp.hashtab))
}

// Builtin is an R language built-in function.
type Builtin struct {
	prim_sexprec
}

// Offset returns the offset into the table of language primitives.
func (v *Builtin) Offset() int32 {
	return int32(v.prim_sxp.offset)
}

// Special is an R language built-in function.
type Special struct {
	prim_sexprec
}

// Offset returns the offset into the table of language primitives.
func (v *Special) Offset() int32 {
	return int32(v.prim_sxp.offset)
}
