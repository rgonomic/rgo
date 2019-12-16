# rgo

![[Paper nautilus](https://archive.org/details/icefalopodiviven00jatt)](Argonauta_argo.png)

`rgo` will be an R/Go integration tool: **WIP**

Install the `rgo` executable.

```
$ go get github.com/rgnonomic/rgo
```

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
$ rgo init
```

Edit the configuration if needed.

```
$ vim rgo.json
```

Generate the wrapper code, documentation and other associated files.

```
$ rgo build example.org/path/to/go_code@vX.Y.Z
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
