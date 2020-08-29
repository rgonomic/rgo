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

// Keep this reasonably in sync with the README.md.
const helpText = `
Workflow

Initialize the go package module.

$ go mod init example.org/path/to/module

Specify dependency versions using the go tool if needed.

$ go get example.org/path/to/dependency@vX.Y.Z

Set up the rgo default configurations and files.

$ rgo init example.org/path/to/pkg

Edit the configuration file, rgo.json, if needed. The rgo.json file
corresponds to the following Go struct.

type Config struct {
	// PkgPath is the package import path for the package
	// to be wrapped by rgo. It depends on the go.mod file
	// at the root of the destination rgo package.
	PkgPath string

	// AllowedFuncs is a pattern matching names of
	// functions that may be wrapped. If AllowedFuncs
	// is empty all wrappable functions are wrapped.
	AllowedFuncs string

	// Words is a set of known words that can be provided
	// to ensure camel-case to snake case breaks words
	// correctly. If words is nil, "NaN" and "NA" are
	// used. Set words to []string{} to provide an empty
	// set of words.
	Words []string

	// LicenseDir is the directory to put license files
	// when more than one license exists.
	LicenseDir string

	// LicensePattern is the pattern for license file
	// names to check. The pattern is used with the
	// case-insensitive flag.
	LicensePattern string
}

Generate the wrapper code, documentation and other associated files.

$ rgo build

Make necessary changes to the DESCRIPTION file.
If the R package is intended to be packaged license information for
dependencies of the Go code will need to be included since Go links
statically and the generated .so lib file will be part of the
distribution. rgo build will collect all the licenses that it finds in
the source module and place them in the LicenseDir directory. You should
remove any that are not relevant to the package you are wrapping.

Then package into a bundle for distribution and build and install the R
package.

$ R CMD INSTALL .


Type mappings

rgo has builtin type mappings between Go and R types. See the project
README for the details of these.


Go struct tags

Go struct tags with the name rgo may be used to change the R value's
name mapping. For example,

type GoType struct {
	Count int ` + "`rgo:\"number\"`" + `
}

will correspond to an R list with a single named element "number".


Multiple return values

Go functions returning multiple values will have these values packaged
into a list with elements named for the return values in the case of Go
functions named returns, or 'r<n>' for unnamed returns where '<n>' is
the index of the return value.


Panics

Go panics are recovered and result in an R error call.


Limitations

R and Go have differences in indexing; R is one-based and Go is
zero-based. This means that care needs to be taken when using indexes
generated in the other environment.

R lacks 64-bit integers, so rgo will refuse to wrap functions that have
64-bit integer inputs or results (int64 and uint64). It also refuses to
wrap function that take or return uintptr values. On Go architectures
with 64-bit int and uint types, results are truncated to 32 bits. This
behaviour will not change until R gets 64-bit integer types.

R Matrix values are not currently handled and will need to be
destructured to a vector and a pair of dimensions (see the example in
examples/cca in the project repo for how to do this).

Currently the extraction of type identities is weaker than it should be.
This will be improved.

Data exchange between R and Go depends on Cgo calls and so is not free.
The exact performance impact depends on the type due to R's baroque type
system and its implementation; briefly though, R vectors that have a
direct correspondence with Go scalar types or slices will perform the
best (integer and int32/uint32, double and float64, complex and
complex128, and raw and int8/uint8). To check the likely performance of
data exchange, look at the generated Go code in the src/rgo directory of
the package you are building. The generated code is intended to be
reasonably human readable.


Input parameter mutation

For types that have direct memory layout equivalents between Go and R
(raw and []int8/[]uint8, integer and []int32/[]uint32, double/float64,
and complex and []complex128) the vector is passed directly to Go. This
means that the Go code can mutate elements. This needs to be considered
when writing Go code that works on slices to avoid unwanted mutation of
R values that are passed to Go. It can also be used for allocation free
work on R vectors. Values passed back to R from Go are copied to satisfy
Go's runtime restrictions on pointer passing.
`
