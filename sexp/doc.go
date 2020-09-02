// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sexp provides an API for interacting with R SEXP values.
// Except where noted, values are allocated by the R runtime.
//
// Building the sexp package requires having CGO_LDFLAGS_ALLOW="-Wl,-Bsymbolic-functions"
// and installations of R and pkg-config.
package sexp

// BUG(kortschak): Currently the API does not provide the capacity to create all SEXP types.
