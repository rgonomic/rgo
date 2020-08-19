// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rgo handles the rgo command line.
// It contains a handler for each of the modes, along with all the flag handling
// and the command line output format.
package rgo

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/rgonomic/rgo/internal/tool"
)

// Application is the main application as passed to tool.Main
// It handles the main command line parsing and dispatch to the sub commands.
type Application struct {
	// Core application flags

	// Embed the basic profiling flags supported by the tool package
	tool.Profile

	// The name of the binary, used in help and telemetry.
	name string

	// The working directory to run commands in.
	wd string

	// The environment variables to use.
	env []string

	// Enable verbose logging
	Verbose bool `flag:"v" help:"verbose output"`
}

// Returns a new Application ready to run.
func New(name, wd string, env []string) *Application {
	if wd == "" {
		wd, _ = os.Getwd()
	}
	return &Application{
		name: name,
		wd:   wd,
		env:  env,
	}
}

// Name implements tool.Application returning the binary name.
func (app *Application) Name() string { return app.name }

// Usage implements tool.Application returning empty extra argument usage.
func (app *Application) Usage() string { return "<command> [command-flags] [command-args]" }

// ShortHelp implements tool.Application returning the main binary help.
func (app *Application) ShortHelp() string {
	return "The rgo source building tools."
}

// DetailedHelp implements tool.Application returning the main binary help.
// This includes the short help for all the sub commands.
func (app *Application) DetailedHelp(f *flag.FlagSet) {
	fmt.Fprint(f.Output(), `
Available commands are:
`)
	for _, c := range app.commands() {
		fmt.Fprintf(f.Output(), "  %s: %v\n", c.Name(), c.ShortHelp())
	}
	fmt.Fprint(f.Output(), `
rgo flags are:
`)
	f.PrintDefaults()
}

// Run takes the args after top level flag processing, and invokes the correct
// sub command as specified by the first argument.
// If no arguments are passed it will invoke the server sub command, as a
// temporary measure for compatibility.
func (app *Application) Run(ctx context.Context, args ...string) error {
	if len(args) == 0 {
		return tool.Run(ctx, &help{}, args)
	}
	command, args := args[0], args[1:]
	for _, c := range app.commands() {
		if c.Name() == command {
			return tool.Run(ctx, c, args)
		}
	}
	return tool.CommandLineErrorf("Unknown command %v", command)
}

// commands returns the set of commands supported by the rgo tool on the
// command line.
// The command is specified by the first non flag argument.
func (app *Application) commands() []tool.Application {
	return []tool.Application{
		&setup{app: app, DryRun: false},
		&build{app: app, DryRun: false},
		&version{app: app},
		&help{},
	}
}
