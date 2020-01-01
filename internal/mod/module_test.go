// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mod

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestModule(t *testing.T) {
	mod, err := ioutil.TempDir("", "rgotest")
	if err != nil {
		t.Fatalf("unexpected error creating mod test directory: %v", err)
	}
	defer os.RemoveAll(mod)
	cmd := exec.Command("go", "mod", "init", "rgotest")
	cmd.Dir = mod
	err = cmd.Run()
	if err != nil {
		t.Fatalf("unexpected error creating mod file: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(mod, "src.go"), []byte("package rgotest"), 0o664)
	if err != nil {
		t.Fatalf("unexpected error creating mod-case source file: %v", err)
	}

	nomod, err := ioutil.TempDir("", "rgotest")
	if err != nil {
		t.Fatalf("unexpected error creating non-mod test directory: %v", err)
	}
	defer os.RemoveAll(nomod)
	err = ioutil.WriteFile(filepath.Join(nomod, "src.go"), []byte("package rgotest"), 0o664)
	if err != nil {
		t.Fatalf("unexpected error creating non-mod-case source file: %v", err)
	}

	oldmode := os.Getenv("GO111MODULE")
	defer os.Setenv("GO111MODULE", oldmode)
	for _, path := range []string{mod, nomod} {
		for _, mode := range []string{"on", "off"} {
			os.Setenv("GO111MODULE", mode)
			got, err := Module(path)
			if path == mod && mode == "on" {
				want := Info{
					Path:  "rgotest",
					Dir:   path,
					GoMod: filepath.Join(path, "go.mod"),
				}
				if got != want {
					t.Errorf("unexpected module info: got:%+v want:%+v", got, want)
				}
			} else {
				var empty Info
				if got != empty {
					t.Errorf("expected empty module info for go.mod!=%t/GO111MODULE=%s", path != mod, mode)
				}
				if mode == "on" {
					// No mod file and GO111MODULE=on outside
					// $GOPATH results in an error:
					//  go: cannot find main module; see 'go help modules'
					continue
				}

			}
			if err != nil {
				t.Errorf("unexpected error getting module information: %v", err)
			}
		}
	}
}

func TestRoot(t *testing.T) {
	mod, err := ioutil.TempDir("", "rgotest")
	if err != nil {
		t.Fatalf("unexpected error creating mod test directory: %v", err)
	}
	defer os.RemoveAll(mod)
	err = os.Mkdir(filepath.Join(mod, "a"), 0o775)
	if err != nil {
		t.Fatalf("unexpected error creating mod test sub-directory: %v", err)
	}
	cmd := exec.Command("go", "mod", "init", "rgotest")
	cmd.Dir = mod
	err = cmd.Run()
	if err != nil {
		t.Fatalf("unexpected error creating mod file: %v", err)
	}

	nomod, err := ioutil.TempDir("", "rgotest")
	if err != nil {
		t.Fatalf("unexpected error creating non-mod test directory: %v", err)
	}
	defer os.RemoveAll(nomod)
	err = os.Mkdir(filepath.Join(nomod, "a"), 0o775)
	if err != nil {
		t.Fatalf("unexpected error creating non-mod test sub-directory: %v", err)
	}

	oldmode := os.Getenv("GO111MODULE")
	defer os.Setenv("GO111MODULE", oldmode)
	for _, path := range []string{mod, nomod} {
		for _, subdir := range []string{"", "a"} {
			for _, mode := range []string{"on", "off"} {
				os.Setenv("GO111MODULE", mode)
				got, ok, err := Root(filepath.Join(path, subdir))
				switch {
				case mode == "on" && path == mod:
					want := path
					if got != want {
						t.Errorf("unexpected root for go.mod=true/GO111MODULE=on in %s: got:%q want:%q",
							filepath.Join(path, subdir), got, want)
					}
					if !ok {
						t.Error("expected ok==true for go.mod=true/GO111MODULE=on")
					}
					if err != nil {
						t.Errorf("unexpected error for go.mod=true/GO111MODULE=on: %v", err)
					}
				case mode == "on" && path == nomod:
					want := filepath.Join(path, subdir)
					if got != want {
						t.Errorf("unexpected root for go.mod=false/GO111MODULE=on in %s: got:%q want:%q",
							filepath.Join(path, subdir), got, want)
					}
					if ok {
						t.Error("expected ok==false for go.mod=false/GO111MODULE=on")
					}
					if err != nil {
						t.Errorf("unexpected error for go.mod=false/GO111MODULE=on: %v", err)
					}
				case mode == "off":
					want := ""
					if got != want {
						t.Errorf("unexpected root for go.mod=%t/GO111MODULE=off in %s: got:%q want:%q",
							path == mod, filepath.Join(path, subdir), got, want)
					}
					if ok {
						t.Errorf("expected ok==false for go.mod=%t/GO111MODULE=off", path == mod)
					}
					if err == nil {
						t.Errorf("expected error for go.mod=%t/GO111MODULE=off: %v", path == mod, err)
					}
				}
			}
		}
	}
}
