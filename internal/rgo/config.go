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

	// LicenseDir is the directory to put license files
	// when more than one license exists.
	LicenseDir string
	// LicensePattern is the pattern for license file
	// names to check. The pattern is used with the
	// case-insensitive flag.
	LicensePattern string
}
