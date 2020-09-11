// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sexp provides an API for interacting with R SEXP values.
// Except where noted, values are allocated by the R runtime.
//
// Here be dragons
//
// The API provided by the sexp package only minimally interacts with
// the R runtime, so deviations from behavior expect by the R runtime
// are possible. This package should be used with caution and client
// packages should ensure correct behavior with adequate testing from
// the R side.
//
// Building the sexp package requires having CGO_LDFLAGS_ALLOW="-Wl,-Bsymbolic-functions"
// and installations of R and pkg-config.
package sexp

// BUG(kortschak): Currently the API does not provide the capacity to create all SEXP types.
