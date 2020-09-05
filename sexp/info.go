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

//go:generate bash -c "go tool cgo -godefs -- $(pkg-config --cflags libR) sxptypes.go | gofmt > cgo_types.go"
//go:generate rm -rf _obj
//go:generate stringer -type=Type

// Type is the SEXPTYPE enum defined in Rinternals.h.
type Type byte
