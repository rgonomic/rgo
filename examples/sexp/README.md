# Direct access to SEXP values

Brave users can directly access the R `SEXP` value by accepting/returning `unsafe.Pointer` parameters (see the [sexp package docs](https://pkg.go.dev/github.com/rgonomic/rgo/sexp?tab=doc) for more details).

```
// PrintSEXP prints the SEXP value passed to it and returns it unaltered.
// If the value is an atomic vector, its contents are printed.
func PrintSEXP(p unsafe.Pointer) unsafe.Pointer {
	sxp := (*sexp.Value)(p)
	info := sxp.Info()
	switch sxp := sxp.Interface().(type) {
	case *sexp.Integer:
		fmt.Printf("%s %d\n", info, sxp.Vector())
	case *sexp.Logical:
		fmt.Printf("%s %d\n", info, sxp.Vector())
	case *sexp.Real:
		fmt.Printf("%s %f\n", info, sxp.Vector())
	case *sexp.Complex:
		fmt.Printf("%s %f\n", info, sxp.Vector())
	case *sexp.String:
		fmt.Printf("%s %s\n", info, sxp.Vector())
	case *sexp.Character:
		fmt.Printf("%s %s\n", info, sxp)
	case *sexp.Raw:
		fmt.Printf("%s %v\n", info, sxp.Bytes())
	default:
		fmt.Println(info)
	}
	return p
}
```

Use of this functionality is beyond the scope of a small example README.

```
> library(sexp)
> sexp::print_sexp(c(1L))
Info{Type:INTSXP, IsScalar:true, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:1, Named:1, Extra:0b0000000000000000} []int32{1}
[1] 1
> sexp::print_sexp("text")
Info{Type:STRSXP, IsScalar:true, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:1, Named:4, Extra:0b0000000000000000} ["text"]
[1] "text"
> sexp::print_sexp(c("more", "text"))
Info{Type:STRSXP, IsScalar:false, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:2, Named:1, Extra:0b0000000000000000} ["more" "text"]
[1] "more" "text"
> sexp::print_sexp(c("even", "more", NA))
Info{Type:STRSXP, IsScalar:false, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:3, Named:1, Extra:0b0000000000000000} ["even" "more" "NA"]
[1] "even" "more" NA    
> sexp::print_sexp(c(TRUE, FALSE, NA))
Info{Type:LGLSXP, IsScalar:false, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:2, Named:1, Extra:0b0000000000000000} []int32{1, 0, -2147483648}
[1]  TRUE FALSE    NA
> sexp::print_sexp(c(1, 2, NA))
Info{Type:REALSXP, IsScalar:false, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:3, Named:1, Extra:0b0000000000000000} []float64{1, 2, NaN}
[1]  1  2 NA
> sexp::print_sexp(c(1L, 2L, NA))
Info{Type:INTSXP, IsScalar:false, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:2, Named:1, Extra:0b0000000000000000} []int32{1, 2, -2147483648}
[1]  1  2 NA
```