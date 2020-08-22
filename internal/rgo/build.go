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
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/rgonomic/rgo/internal/codegen"
	"github.com/rgonomic/rgo/internal/mod"
	"github.com/rgonomic/rgo/internal/pkg"
	"github.com/rgonomic/rgo/internal/vfs/osfs"
	"github.com/rgonomic/rgo/internal/vfs/txtar"
)

// TODO(kortschak): Allow `rgo build [root]`.

// build implements the build command.
type build struct {
	app *Application

	Config

	DryRun bool `flag:"dry-run" help:"output generated code to stdout as a txtar"`
}

func (b *build) Name() string      { return "build" }
func (b *build) Usage() string     { return "" }
func (b *build) ShortHelp() string { return "runs the rgo build command" }
func (b *build) DetailedHelp(f *flag.FlagSet) {
	fmt.Fprint(f.Output(), ``)
	f.PrintDefaults()
}

// Run runs the rgo build command.
func (b *build) Run(ctx context.Context, args ...string) error {
	root, ok, err := mod.Root(b.app.wd)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf(`%s is not within a go module: run "go mod init"`, b.app.wd)
	}

	cfg, err := ioutil.ReadFile(filepath.Join(root, "rgo.json"))
	if err != nil {
		return fmt.Errorf("failed to read config: %w. run rgo init", err)
	}
	err = json.Unmarshal(cfg, &b.Config)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	candidate, err := regexp.Compile(fmt.Sprintf("(?i:%s)", b.LicensePattern))
	if err != nil {
		return fmt.Errorf("failed to parse license name pattern: %w", err)
	}

	info, err := pkg.Analyse(b.PkgPath, b.AllowedFuncs, b.app.Verbose)
	if err != nil {
		return fmt.Errorf("load error: %w", err)
	}
	if len(info.Funcs) == 0 {
		log.Println("no functions to wrap")
		return nil
	}

	if b.app.Verbose {
		log.Println("need C.SEXP->Go for:")
		for _, bt := range info.Unpackers {
			log.Printf(" %s: %s\n", pkg.Mangle(bt), bt)
		}
		log.Println("need Go->C.SEXP for:")
		for _, bt := range info.Packers {
			log.Printf(" %s: %s\n", pkg.Mangle(bt), bt)
		}
	}

	var dst codegen.FileSystem
	if b.DryRun {
		dst = &txtar.FileSystem{}
	} else {
		dst = osfs.FileSystem{}
	}

	err = codegen.Licenses(dst, b.LicenseDir, info, candidate, b.app.Verbose)
	if err != nil {
		return fmt.Errorf("failed to get license information: %w", err)
	}

	words := []string{"NaN", "NA"}
	templates := map[string]*template.Template{
		"NAMESPACE":     codegen.NamespaceTemplate(words),
		"R/%s.R":        codegen.RCallTemplate(words),
		"src/rgo/%s.c":  codegen.CFuncTemplate(words),
		"src/rgo/%s.go": codegen.GoFuncTemplate(),
		"src/Makevars":  codegen.MakevarsTemplate(),
	}
	for path, tmpl := range templates {
		err = codegen.Render(dst, path, tmpl, info)
		if err != nil {
			return fmt.Errorf("failed to render %s: %w", tmpl.Name(), err)
		}
	}

	err = dst.Flush()
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	return nil
}
