// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rgo

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
)

// version implements the version command.
type version struct {
	app *Application
}

func (v *version) Name() string      { return "version" }
func (v *version) Usage() string     { return "" }
func (v *version) ShortHelp() string { return "print the rgo version information" }
func (v *version) DetailedHelp(f *flag.FlagSet) {
	fmt.Fprint(f.Output(), ``)
	f.PrintDefaults()
}

// Run prints rgo version information.
func (v *version) Run(ctx context.Context, args ...string) error {
	printBuildInfo(os.Stdout, v.app.Verbose)
	return nil
}

func printBuildInfo(w io.Writer, verbose bool) {
	if info, ok := debug.ReadBuildInfo(); ok {
		fmt.Fprintf(w, "%v %v\n", info.Path, info.Main.Version)
		if verbose {
			for _, dep := range info.Deps {
				printModuleInfo(w, dep)
			}
		}
	} else {
		fmt.Fprintf(w, "version unknown, built in $GOPATH mode\n")
	}
}

func printModuleInfo(w io.Writer, m *debug.Module) {
	fmt.Fprintf(w, "    %s@%s", m.Path, m.Version)
	if m.Sum != "" {
		fmt.Fprintf(w, " %s", m.Sum)
	}
	if m.Replace != nil {
		fmt.Fprintf(w, " => %v", m.Replace.Path)
	}
	fmt.Fprintf(w, "\n")
}

// help implements the help command.
type help struct{}

func (*help) Name() string      { return "help" }
func (*help) Usage() string     { return "" }
func (*help) ShortHelp() string { return "output rgo help information" }
func (*help) DetailedHelp(f *flag.FlagSet) {
	fmt.Fprint(f.Output(), ``)
	f.PrintDefaults()
}

// Run outputs the help text.
func (*help) Run(ctx context.Context, args ...string) error {
	fmt.Fprintf(os.Stdout, "%s", helpText)
	return nil
}

const helpText = `
TODO(kortschak): Add text.
`
