// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mod

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Module returns module information for the given directory path.
// If the path is in the standard library, the Info Path field will
// be GOROOT and Version will be the Go version.
func Module(path string) (Info, error) {
	args := []string{"list", "-json"}
	if path == "." {
		args = append(args, "-m")
	} else {
		args = append(args, path)
	}
	cmd := exec.Command("go", args...)
	var buf, errbuf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		stderr := strings.TrimSpace(errbuf.String())
		if stderr == "" {
			return Info{}, err
		}
		return Info{}, fmt.Errorf("%s: %w", stderr, err)
	}
	var m struct {
		Module   Info
		Standard bool
	}
	err = json.Unmarshal(buf.Bytes(), &m)
	if err != nil {
		return Info{}, err
	}
	if !m.Standard {
		return m.Module, nil
	}

	buf.Reset()
	errbuf.Reset()
	cmd = exec.Command("go", "env", "GOROOT")
	cmd.Stdout = &buf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		stderr := strings.TrimSpace(errbuf.String())
		if stderr == "" {
			return Info{}, err
		}
		return Info{}, fmt.Errorf("%s: %w", stderr, err)
	}
	root := strings.TrimSpace(buf.String())
	info := Info{
		Path: filepath.Join(root, "src", path),
		Dir:  root,
	}

	buf.Reset()
	errbuf.Reset()
	cmd = exec.Command("go", "version")
	cmd.Stdout = &buf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		stderr := strings.TrimSpace(errbuf.String())
		if stderr == "" {
			return Info{}, err
		}
		return Info{}, fmt.Errorf("%s: %w", stderr, err)
	}
	version := strings.Fields(strings.TrimPrefix(buf.String(), "go version go"))[0]
	info.Version = "v" + version

	return info, nil
}

// Info holds module information. It is a subset of the information
// returned by "go list -json" in the Module struct. See go help list
// for details.
type Info struct {
	Path    string // module path
	Version string // module version
	Dir     string // directory holding files for this module, if any
	GoMod   string // path to go.mod file for this module, if any
}

// Root returns the root directory of the module in the given path and
// whether a go.mod file can be found. It returns an error if the go
// tool is not running is module-aware mode.
func Root(path string) (root string, ok bool, err error) {
	cmd := exec.Command("go", "env", "GOMOD")
	cmd.Dir = path
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err = cmd.Run()
	if err != nil {
		return "", false, err
	}
	gomod := strings.TrimSpace(buf.String())
	if gomod == "" {
		return "", false, errors.New("go tool not running in module mode")
	}
	if gomod == os.DevNull {
		return path, false, nil
	}
	return filepath.Dir(gomod), true, nil
}
