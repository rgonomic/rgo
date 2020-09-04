// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sexp

import (
	"fmt"
	"unsafe"
)

//go:generate go run generate_types.go

// Pointer returns an unsafe pointer to the SEXP value.
func (v *sexprec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

// Pointer returns an unsafe pointer to the SEXP value.
func (v *vector_sexprec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

// Pointer returns an unsafe pointer to the SEXP value.
func (v *list_sexprec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

// Pointer returns an unsafe pointer to the SEXP value.
func (v *env_sexprec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

// Pointer returns an unsafe pointer to the SEXP value.
func (v *prom_sexprec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

// Pointer returns an unsafe pointer to the SEXP value.
func (v *prim_sexprec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

// Pointer returns an unsafe pointer to the SEXP value.
func (v *sym_sexprec) Pointer() unsafe.Pointer {
	return unsafe.Pointer(v)
}

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
