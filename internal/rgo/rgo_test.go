// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rgo

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/diff"
)

var regenerate = flag.Bool("regen", false, "regenerate golden data from current state")

func TestRgo(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "rgo-*")
	if err != nil {
		t.Fatalf("failed to make temporary build directory: %v", err)
	}
	rgo := filepath.Join(tmpdir, "rgo")
	cmd := exec.Command("go", "build", "-o", rgo, "github.com/rgonomic/rgo")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build rgo: %v\n%s", err, out)
	}
	t.Cleanup(func() {
		err := os.RemoveAll(tmpdir)
		if err != nil {
			t.Fatalf("failed to clean up build directory: %v", err)
		}
	})

	dir, err := os.Open("testdata")
	if err != nil {
		t.Fatalf("failed open testdata directory: %v", err)
	}
	infos, err := dir.Readdir(-1)
	if err != nil {
		t.Fatalf("failed open get package names: %v", err)
	}
	for _, fi := range infos {
		if !fi.IsDir() {
			continue
		}
		pkg := fi.Name()
		t.Run("init:"+pkg, func(t *testing.T) {
			cmd := exec.Command(rgo, "init", fmt.Sprintf("-dry-run=%t", !*regenerate))
			cmd.Dir = filepath.Join("testdata", pkg)
			var buf bytes.Buffer
			cmd.Stdout = &buf
			err := cmd.Run()
			if err != nil {
				t.Fatalf("failed to run rgo init: %v", err)
			}

			if *regenerate {
				return
			}

			got := buf.Bytes()
			want, err := ioutil.ReadFile(filepath.Join("testdata", pkg, "rgo.json"))
			if err != nil {
				t.Fatalf("failed to read golden data: %v", err)
			}

			if !bytes.Equal(got, want) {
				t.Errorf("unexpected rgo.json:\ngot:\n:%s\nwant:\n%s", got, want)
			}
		})

		t.Run("build:"+pkg, func(t *testing.T) {
			if strings.Contains(pkg, "uintptr") {
				t.Skipf("skipping unhandled type %q", "uintptr")
			}
			cmd := exec.Command(rgo, "build", "-dry-run")
			cmd.Dir = filepath.Join("testdata", pkg)
			var buf bytes.Buffer
			cmd.Stdout = &buf
			err := cmd.Run()
			if err != nil {
				t.Fatalf("failed to run rgo build: %v", err)
			}

			got := buf.Bytes()
			golden := filepath.Join("testdata", pkg, "golden.txtar")
			if *regenerate {
				err := ioutil.WriteFile(golden, got, 0o664)
				if err != nil {
					t.Fatalf("failed to write golden data: %v", err)
				}
				return
			}

			want, err := ioutil.ReadFile(golden)
			if err != nil {
				t.Fatalf("failed to read golden data: %v", err)
			}

			if !bytes.Equal(got, want) {
				var buf bytes.Buffer
				err := diff.Text("got", "want", got, want, &buf)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				t.Errorf("unexpected generated code:\n%s", &buf)
			}
		})
	}
}
