# rgo

![[Paper nautilus](https://archive.org/details/icefalopodiviven00jatt)](Argonauta_argo.png)

`rgo` is an R/Go integration tool — it is only lightly tested so far and should not be used in production.

Install the `rgo` executable.

```
$ go get github.com/rgnonomic/rgo
```

## Workflow

With the `rgo` executable installed, the general work-flow will then proceed something like this.

Initialize the go package module.

```
$ go mod init example.org/path/to/module
```

Specify dependency versions using the go tool if needed.

```
$ go get example.org/path/to/dependency@vX.Y.Z
```

Set up the rgo default configurations and files.

```
$ rgo init example.org/path/to/pkg
```

Edit the configuration file, `rgo.json`, if needed. The `rgo.json` file corresponds to the [following Go struct](https://pkg.go.dev/github.com/rgonomic/rgo/internal/rgo?tab=doc#Config).
```
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
```

```
$ vim rgo.json
```

Generate the wrapper code, documentation and other associated files.

```
$ rgo build
```

Make necessary changes to the DESCRIPTION file.
If the R package is intended to be packaged license information for dependencies of the Go code will need to be included since Go links statically and the generated .so lib file will be part of the distribution. `rgo build` will collect all the licenses that it finds in the source module and place them in the LicenseDir directory. You should remove any that are not relevant to the package you are wrapping.

```
$ vim DESCRIPTION
```

Then package into a bundle for distribution...

```
$ git ...
```

... build and install the R package.

```
$ R CMD INSTALL .
```

## Type mappings

`rgo` has builtin type mappings between Go and R types. These are described here.

| R                               | Go                                                                                 |
|---------------------------------|------------------------------------------------------------------------------------|
| scalar `integer`                | `int`, `int8`, `int16`, `int32`/`rune`, `uint`, `uint8`/`byte`, `uint16`, `uint32` |
| `integer` vector                | `[]int`, `[]int16`, `[]int32`/`[]rune`, `[]uint`, `[]uint16`, `[]uint32`           |
| fixed length `integer` vector   | `[n]int`, `[n]int16`, `[n]int32`/`[n]rune`, `[n]uint`, `[n]uint16`, `[n]uint32`    |
| scalar `double`                 | `float32`, `float64`                                                               |
| `double` vector                 | `[]float32`, `[]float64`                                                           |
| fixed length `double` vector    | `[n]float32`, `[n]float64`                                                         |
| scalar `complex`                | `complex64`, `complex128`                                                          |
| `complex` vector                | `[]complex64`, `[]complex128`                                                      |
| fixed length `complex` vector   | `[n]complex64`, `[n]complex128`                                                    |
| scalar `logical`                | `bool`                                                                             |
| `logical` vector                | `[]bool`                                                                           |
| fixed length `logical` vector   | `[n]bool`                                                                          |
| scalar `character`              | `string` (and `error` in returned values)                                          |
| `character` vector              | `[]string` (and `[]error` in returned values)                                      |
| fixed length `character` vector | `[n]string` (and `[n]error` in returned values)                                    |
| unnamed `list`                  | `[]C`                                                                              |
| fixed length unnamed `list`     | `[n]C`                                                                             |
| named `vector`                  | `map[string]A`                                                                     |
| named `list`                    | `map[string]C`                                                                     |
| `list`                          | `struct{...}`                                                                      |
| `raw`                           | `[]int8`, `[]uint8`/`[]byte`                                                       |
| fixed length `raw`              | `[n]int8`, `[n]uint8`/`[n]byte`                                                    |
| internal `SEXP` value           | `unsafe.Pointer`                                                                   |

The Go `A` types correspond to R `atomic` types.

> `int`, `int8`, `int16`, `int32`/`rune`, `uint`, `uint8`/`byte`, `uint16`, `uint32`, `float32`, `float64`, `complex64`, `complex128`, `bool`, `string` and `error`

The Go `C` types correspond all other supported types.

> `[]T`, `[n]T`, `map[string]T` and `struct{...}` where `T` is any type.


Pointer types are also handled. Currently pointers are indirected so that mutations to pointees do not propagate between the Go and R environments. This behaviour may change for pointers being passed to Go from R.


### Go struct tags

Go struct tags with the name `rgo` may be used to change the R value's name mapping. For example,

```
type GoType struct {
	Count int `rgo:"number"`
}
```

will correspond to an R `list` with a single named element `number`.


### Multiple return values

Go functions returning multiple values will have these values packaged into a list with elements named for the return values in the case of Go functions named returns, or `r<n>` for unnamed returns where `<n>` is the index of the return value.


## Panics

Go panics are recovered and result in an R error call.


## Limitations

R and Go have differences in indexing; R is one-based and Go is zero-based. This means that care needs to be taken when using indexes generated in the other environment.

R lacks 64-bit integers, so `rgo` will refuse to wrap functions that have 64-bit integer inputs or results (`int64` and `uint64`). It also refuses to wrap function that take or return `uintptr` values. On Go architectures with 64-bit `int` and `uint` types, results are truncated to 32 bits. This behaviour will not change until R gets 64-bit integer types.

R Matrix values are not currently handled and will need to be destructured to a vector and a pair of dimensions (see the [matrix example](examples/cca) for how to do this).

Currently the extraction of type identities is weaker than it should be. This will be improved.

Data exchange between R and Go depends on Cgo calls and so is not free. The exact performance impact depends on the type due to R's baroque type system and its implementation; briefly though, R vectors that have a direct correspondence with Go scalar types or slices will perform the best (`integer` and `int32`/`uint32`, `double` and `float64`, `complex` and `complex128`, and `raw` and `int8`/`uint8`). To check the likely performance of data exchange, look at the generated Go code in the `src/rgo` directory of the package you are building. The generated code is intended to be reasonably human readable.

R does not know how to unload so libraries, so once loaded it is there until the session is restarted.

## Input parameter mutation

For types that have direct memory layout equivalents between Go and R (`raw` and`[]int8`/`[]uint8`, `integer` and `[]int32`/`[]uint32`, `double` and `[]float64`, and `complex` and`[]complex128`) the vector is passed directly to Go. This means that the Go code can mutate elements. This needs to be considered when writing Go code that works on slices to avoid unwanted mutation of R values that are passed to Go. It can also be used for allocation free work on R vectors. Values passed back to R from Go are copied to satisfy Go's runtime restrictions on pointer passing.