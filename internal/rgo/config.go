// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rgo

// Config is an rgo build config.
type Config struct {
	// PkgPath is the package import path for the package
	// to be wrapped by rgo. It depends on the go.mod file
	// at the root of the destination rgo package.
	PkgPath string

	// AllowedFuncs is a pattern matching names of
	// functions that may be wrapped. If AllowedFuncs
	// is empty all wrappable functions are wrapped.
	AllowedFuncs string

	// Words is a set of known words that can be provided
	// to ensure camel-case to snake case breaks words
	// correctly. If words is nil, "NaN" and "NA" are
	// used. Set words to []string{} to provide an empty
	// set of words.
	Words []string

	// LicenseDir is the directory to put license files
	// when more than one license exists.
	LicenseDir string

	// LicensePattern is the pattern for license file
	// names to check. The pattern is used with the
	// case-insensitive flag.
	LicensePattern string
}
