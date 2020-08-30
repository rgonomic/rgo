# Direct access to SEXP values

Brave users can directly access the R `SEXP` value by accepting/returning `unsafe.Pointer` parameters (see [sexp.go](sexp.go) for more details).

```
func PrintSEXPHeader(sexp unsafe.Pointer) unsafe.Pointer {
	p := (*SEXP)(sexp)
	fmt.Printf("%#v\n", p.Header)
	return sexp
}
```

Use of this functionality is beyond the scope of a small example README.

```
> library(sexp)
> sexp::print_sexp_header(1L)
R.SEXP{Type:INTSXP, IsScalar:true, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:1, Named:4, Extra:0b0000000000000000}
[1] 1
> sexp::print_sexp_header("text")
R.SEXP{Type:STRSXP, IsScalar:true, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:1, Named:4, Extra:0b0000000000000000}
[1] "text"
> sexp::print_sexp_header(c("more", "text"))
R.SEXP{Type:STRSXP, IsScalar:false, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:2, Named:1, Extra:0b0000000000000000}
[1] "more" "text"
```