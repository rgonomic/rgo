// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package osfs provides a basic write-only file system interface to the
// underlying OS file system.
package osfs

import (
	"io"
	"os"
	"path/filepath"
)

// FileSystem is a handle to the OS file system.
type FileSystem struct{}

// Open returns a file for writing, creating any necessary directories.
// The file must be closed after use.
func (FileSystem) Open(path string) (io.WriteCloser, error) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0o775)
	if err != nil {
		return nil, err
	}
	return os.Create(path)
}

// Flush is a no-op.
func (fs FileSystem) Flush() error { return nil }
