// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rgo

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rgonomic/rgo/internal/mod"
)

// setup implements the setup command.
type setup struct {
	app *Application

	DryRun bool `flag:"dry-run" help:"output config to stdout"`
}

func (s *setup) Name() string      { return "init" }
func (s *setup) Usage() string     { return "<package path> (defaults to current directory)" }
func (s *setup) ShortHelp() string { return "runs the rgo init command" }
func (s *setup) DetailedHelp(f *flag.FlagSet) {
	fmt.Fprint(f.Output(), ``)
	f.PrintDefaults()
}

// Run runs the rgo setup command.
func (s *setup) Run(ctx context.Context, args ...string) error {
	root := s.app.wd
	var path string
	switch len(args) {
	case 0:
		path = "."
	case 1:
		path = args[0]
	case 2:
		root = args[0]
		path = args[1]
	default:
		return fmt.Errorf("too many arguments")
	}

	root, ok, err := mod.Root(root)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf(`%s is not within a go module: run "go mod init"`, root)
	}
	err = os.Chdir(root)
	if err != nil {
		return fmt.Errorf("cannot change directory to module root: %w", err)
	}
	m, err := mod.Module(path)
	if err != nil {
		return err
	}
	if m.Dir == root {
		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, abs)
		if err != nil {
			return err
		}
		path = filepath.Clean(filepath.Join(m.Path, rel))
	}

	cfg := Config{
		PkgPath:        filepath.ToSlash(path),
		LicenseDir:     "LICENSE",
		LicensePattern: `(LICEN[SC]E|COPYING)(\.(txt|md))?$`,
	}
	b, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	var w io.Writer = os.Stdout
	if !s.DryRun {
		f, err := os.Create(filepath.Join(root, "rgo.json"))
		if err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
		defer f.Close()
		w = f
	}
	fmt.Fprintf(w, "%s\n", b)
	return nil
}
