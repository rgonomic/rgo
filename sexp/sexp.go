// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sexp

import (
	"fmt"
	"math"
	"unsafe"
)

// IsRealNA returns whether f is an NA value.
func IsRealNA(f float64) bool {
	payload, ok := naNPayload(f)
	return ok && payload == 1954
}

// RealNA returns the R real NA value.
func RealNA() float64 {
	return naNWith(1954)
}

// IsComplexNA returns whether c is an NA value. This is true if either
// of the parts of c is NA.
func IsComplexNA(c complex128) bool {
	return IsRealNA(real(c)) || IsRealNA(imag(c))
}

// ComplexNA returns the R complex NA value. The value returned is NA+0i.
func ComplexNA() complex128 {
	return complex(RealNA(), 0)
}

const (
	nanBits = 0x7ff8000000000000
	nanMask = 0xfff8000000000000
)

// naNWith returns an IEEE 754 "quiet not-a-number" value with the
// payload specified in the low 51 bits of payload.
// The NaN returned by math.NaN has a bit pattern equal to naNWith(1).
func naNWith(payload uint64) float64 {
	return math.Float64frombits(nanBits | (payload &^ nanMask))
}

// naNPayload returns the lowest 51 bits payload of an IEEE 754 "quiet
// not-a-number". For values of f other than quiet-NaN, naNPayload
// returns zero and false.
func naNPayload(f float64) (payload uint64, ok bool) {
	b := math.Float64bits(f)
	if b&nanBits != nanBits {
		return 0, false
	}
	return b &^ nanMask, true
}

// IsIntegerNA returns whether i is an NA value.
func IsIntegerNA(i int32) bool {
	return i == r_NaInt
}

// IntegerNA returns the R integer NA value.
func IntegerNA() int32 {
	return r_NaInt
}

// IsLogicalNA returns whether b is an NA value.
func IsLogicalNA(b int32) bool {
	return b == r_NaInt
}

// LogicalNA returns the R logical NA value.
func LogicalNA() int32 {
	return r_NaInt
}

//go:generate go run generate_types.go

// Value is an SEXP value.
type Value struct {
	sexprec
}

// Info returns the information field of the SEXP value.
func (v *Value) Info() Info {
	if v == nil {
		return NilValue.Info()
	}
	return *(*Info)(unsafe.Pointer(&v.sxpinfo))
}

// Value returns the generic state of the SEXP value.
func (v *Value) Value() *Value {
	return v
}

// Attributes returns the attributes of the SEXP value.
func (v *Value) Attributes() *List {
	if v == nil {
		return nil
	}
	attr := (*List)(unsafe.Pointer(v.attrib))
	if attr.Value().IsNull() {
		return nil
	}
	return attr
}

// Pointer returns an unsafe pointer to the SEXP value.
func (v *Value) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

// Export returns an unsafe pointer to the SEXP value first converting
// a nil value to the R NilValue.
func (v *Value) Export() unsafe.Pointer {
	if v == nil {
		return NilValue.Pointer()
	}
	return unsafe.Pointer(v)
}

// IsNull returned whether the Value is the R NULL value.
func (v *Value) IsNull() bool {
	return v.Info().Type() == NILSXP
}

// Interface returns a Go value corresponding to the SEXP type specified
// in the SEXP info field. If the receiver is nil, the R NilValue will be
// returned.
func (v *Value) Interface() interface{} {
	if v == nil {
		return NilValue
	}
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

// Names returns the names of the Value.
func (v *Value) Names() *String {
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
	if attr == nil {
		return nil
	}
	names := attr.Get([]byte("names"))
	found, ok := names.Interface().(*List)
	if !ok {
		return nil
	}
	return found.Head().Interface().(*String)
}
