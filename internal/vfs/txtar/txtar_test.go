// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package txtar

import (
	"bytes"
	"testing"
)

var fileSystemTests = []struct {
	data map[string][]string
	want string
}{
	{
		data: map[string][]string{
			"a": {"A"},
		},
		want: `-- a --
A
`,
	},
	{
		data: map[string][]string{
			"b": {"B"},
			"a": {"A"},
		},
		want: `-- a --
A
-- b --
B
`,
	},
	{
		data: map[string][]string{
			"b": {"B1", "B2"},
			"a": {"A1", "A2"},
		},
		want: `-- a --
A1A2
-- b --
B1B2
`,
	},
	{
		data: map[string][]string{
			"path/to/b": {"B1", "B2"},
			"path/to/a": {"A1", "A2"},
		},
		want: `-- path/to/a --
A1A2
-- path/to/b --
B1B2
`,
	},
}

func TestFileSystem(t *testing.T) {
	for _, test := range fileSystemTests {
		var buf bytes.Buffer
		fs := &FileSystem{Output: &buf}
		for path, writes := range test.data {
			w, err := fs.Open(path)
			for _, d := range writes {
				n, err := w.Write([]byte(d))
				if err != nil {
					t.Errorf("unexpected write error: %v", err)
				}
				if n != len(d) {
					t.Errorf("unexpected number of bytes written for %q: got: %d want: %d", d, n, len(d))
				}
			}
			err = w.Close()
			if err != nil {
				t.Errorf("unexpected close error: %v", err)
			}
		}
		err := fs.Flush()
		if err != nil {
			t.Errorf("unexpected flush error: %v", err)
		}
		got := buf.String()
		if got != test.want {
			t.Errorf("unexpected result:\ngot:\n%s\nwant:\n%s", got, test.want)
		}
	}
}
