// Package osfs provides a basic write-only file system interface to the
// underlying OS file system.
package osfs

import (
	"io"
	"os"
)

// FileSystem is a handle to the OS file system.
type FileSystem struct{}

// Open returns a file for writing. The file must be closed after use.
func (FileSystem) Open(path string) (io.WriteCloser, error) {
	return os.Open(path)
}

// Flush is a no-op.
func (fs *FileSystem) Flush() error { return nil }
