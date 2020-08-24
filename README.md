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

Edit the configuration if needed.

```
$ vim rgo.json
```

Generate the wrapper code, documentation and other associated files.

```
$ rgo build
```

Make necessary changes to the DESCRIPTION file.
If the R package is intended to be packaged license information for dependencies of the Go code will need to be included since Go links statically and the generated .so lib file will be part of the distribution.

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

| R                               | Go                                                                                         |
|---------------------------------|--------------------------------------------------------------------------------------------|
| scalar `integer`                | `int`, `int8`, `int16`, `int32`/`rune`, `uint`, `uint8`/`byte`, `uint16`, `uint32`         |
| `integer` vector                | `[]int`, `[]int8`, `[]int16`, `[]int32`/`[]rune`, `[]uint`, `[]uint16`, `[]uint32`         |
| fixed length `integer` vector   | `[n]int`, `[n]int8`, `[n]int16`, `[n]int32`/`[n]rune`, `[n]uint`, `[n]uint16`, `[n]uint32` |
| scalar `double`                 | `float32`, `float64`                                                                       |
| `double` vector                 | `[]float32`, `[]float64`                                                                   |
| fixed length `double` vector    | `[n]float32`, `[n]float64`                                                                 |
| scalar `complex`                | `complex64`, `complex128`                                                                  |
| `complex` vector                | `[]complex64`, `[]complex128`                                                              |
| fixed length `complex` vector   | `[n]complex64`, `[n]complex128`                                                            |
| scalar `logical`                | `bool`                                                                                     |
| `logical` vector                | `[]bool`                                                                                   |
| fixed length `logical` vector   | `[n]bool`                                                                                  |
| scalar `character`              | `string` (and `error` in returned values)                                                  |
| `character` vector              | `[]string` (and `[]error` in returned values)                                              |
| fixed length `character` vector | `[n]string` (and `[n]error` in returned values)                                            |
| named `vector`                  | `map[string]T`                                                                             |
| `list`                          | `struct{...}`                                                                              |
| `raw`                           | `[]uint8`/`[]byte`                                                                         |
| fixed length `raw`              | `[n]uint8`/`[n]byte`                                                                       |


Pointer types are also handled. Currently pointers are indirected so that mutations to pointees do not propagate between the Go and R environments. This behaviour may change for pointers being passed to Go from R.


### Go struct tags

Go struct tags with the name `rgo` may be used to change the R `list` name mapping. For example,

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

R Matrix values are not currently handled and will need to be destructured to a vector and a pair of dimensions.

Currently the extraction of type identities is weaker than it should be. This will be improved.


## Input parameter mutation

For types that have direct memory layout equivalents between Go and R (`raw`/`[]byte`, `integer`/`[]int32`, `double`/`float64` and `complex`/`[]complex128`) the vector is passed directly to Go. This means that the Go code can mutate elements. This needs to be considered when writing Go code that works on slices to avoid unwanted mutation of R values that are passed to Go. It can also be used for allocation free work on R vectors. Values passed back to R from Go are copied to satisfy Go's runtime restrictions on pointer passing.