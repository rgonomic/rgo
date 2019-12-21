// Package txtar provides a simple write-only filesystem based on the txtar archive format.
package txtar

import (
	"bytes"
	"io"
	"os"
	"sort"

	"golang.org/x/tools/txtar"
)

// FileSystem is a write-only filesystem based on the txtar text archive format
// provided by golang.org/x/tools/txtar.
type FileSystem struct {
	// Output is the destination of the filesystem
	// when the Flush method is called. If Output
	// is nil, data is written to os.Stdout.
	Output io.Writer
	ar     txtar.Archive
	buf    bytes.Buffer
}

// Open returns a new io.WriteCloser file whose writes will added to the
// txtar output of the Filesystem with the given path. Writes to the
// file are buffered until the Flush method is called. The file must
// be closed before calling Flush.
func (fs *FileSystem) Open(path string) (io.WriteCloser, error) {
	return file{path: path, buf: &fs.buf, owner: fs}, nil
}

func (fs *FileSystem) add(path string, data []byte) {
	fs.ar.Files = append(fs.ar.Files, txtar.File{
		Name: path,
		Data: data,
	})
}

// Flush writes all closed io.Writers to the FileSystem's Output
// or os.Stdout if Output is nil. The files in the txtar are sorted
// lexically by name.
func (fs *FileSystem) Flush() error {
	sort.Sort(byPath(fs.ar.Files))
	var dst io.Writer = os.Stdout
	if fs.Output != nil {
		dst = fs.Output
	}
	_, err := dst.Write(txtar.Format(&fs.ar))
	return err
}

type file struct {
	path  string
	buf   *bytes.Buffer
	owner *FileSystem
}

func (f file) Write(b []byte) (int, error) {
	return f.buf.Write(b)
}

func (f file) Close() error {
	b := make([]byte, f.buf.Len())
	copy(b, f.buf.Bytes())
	f.owner.add(f.path, b)
	f.buf.Reset()
	return nil
}

type byPath []txtar.File

func (f byPath) Len() int           { return len(f) }
func (f byPath) Less(i, j int) bool { return f[i].Name < f[j].Name }
func (f byPath) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
