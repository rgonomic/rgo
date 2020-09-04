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
		fmt.Printf("%s values:%#v", info, sxp.Vector())
	case *sexp.Logical:
		fmt.Printf("%s values:%#v", info, sxp.Vector())
	case *sexp.Real:
		fmt.Printf("%s values:%#v", info, sxp.Vector())
	case *sexp.Complex:
		fmt.Printf("%s values:%#v", info, sxp.Vector())
	case *sexp.String:
		fmt.Printf("%s values:%q", info, sxp.Vector())
	case *sexp.Character:
		fmt.Printf("%s values:%s", info, sxp)
	case *sexp.Raw:
		fmt.Printf("%s values:%#v", info, sxp.Bytes())
	default:
		fmt.Print(info)
	}
	fmt.Printf(" names:%v\n", sxp.Value().Names().Vector())
	return p
}
```

Use of this functionality is beyond the scope of a small example README.

```
> library(sexp)
> sexp::print_sexp(c(1L))
Info{Type:INTSXP, IsScalar:true, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:1, Named:1, Extra:0b0000000000000000} values:[]int32{1} names:[]
[1] 1
> sexp::print_sexp("text")
Info{Type:STRSXP, IsScalar:true, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:1, Named:4, Extra:0b0000000000000000} values:["text"] names:[]
[1] "text"
> sexp::print_sexp(c("more", "text"))
Info{Type:STRSXP, IsScalar:false, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:2, Named:1, Extra:0b0000000000000000} values:["more" "text"] names:[]
[1] "more" "text"
> sexp::print_sexp(c("even", "more", NA))
Info{Type:STRSXP, IsScalar:false, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:3, Named:1, Extra:0b0000000000000000} values:["even" "more" "NA"] names:[]
[1] "even" "more" NA    
> sexp::print_sexp(c(TRUE, FALSE, NA))
Info{Type:LGLSXP, IsScalar:false, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:2, Named:1, Extra:0b0000000000000000} values:[]int32{1, 0, -2147483648} names:[]
[1]  TRUE FALSE    NA
> sexp::print_sexp(c(1, 2, NA))
Info{Type:REALSXP, IsScalar:false, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:3, Named:1, Extra:0b0000000000000000} values:[]float64{1, 2, NaN} names:[]
[1]  1  2 NA
> sexp::print_sexp(c(1L, 2L, NA))
Info{Type:INTSXP, IsScalar:false, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:2, Named:1, Extra:0b0000000000000000} values:[]int32{1, 2, -2147483648} names:[]
[1]  1  2 NA
> sexp::print_sexp(c(a=1, b=2))
Info{Type:REALSXP, IsScalar:false, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:2, Named:1, Extra:0b0000000000000000} values:[]float64{1, 2} names:[a b]
a b 
1 2 
> sexp::print_sexp(list(a="text", b=2))
Info{Type:VECSXP, IsScalar:false, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:2, Named:1, Extra:0b0000000000000000} names:[a b]
$a
[1] "text"

$b
[1] 2

> sexp::print_sexp(pairlist(a="text", b=2))
Info{Type:LISTSXP, IsScalar:false, IsObject:false, IsAltRep:false, GP:0b0000000000000000, IsInUse:false, Debug:false, Trace:false, Spare:false, GCGeneration:0, GCClass:0, Named:1, Extra:0b0000000000000000} names:[a b]
$a
[1] "text"

$b
[1] 2
```

New SEXP values can be created with the sexp package, for example, a vector of strings.

```
// Gophers returns n gophers.
func Gophers(n int) unsafe.Pointer {
	c := sexp.NewString(n).Protect()
	defer c.Unprotect()
	vec := c.Vector()
	for i := range vec {
		vec[i] = sexp.NewCharacter(fmt.Sprintf("Gopher %d", i+1))
	}
	return c.Pointer()
}
```

```
> sexp::gophers(1L)
[1] "Gopher 1"
> sexp::gophers(100L)
  [1] "Gopher 1"   "Gopher 2"   "Gopher 3"   "Gopher 4"   "Gopher 5"  
  [6] "Gopher 6"   "Gopher 7"   "Gopher 8"   "Gopher 9"   "Gopher 10" 
 [11] "Gopher 11"  "Gopher 12"  "Gopher 13"  "Gopher 14"  "Gopher 15" 
 [16] "Gopher 16"  "Gopher 17"  "Gopher 18"  "Gopher 19"  "Gopher 20" 
 [21] "Gopher 21"  "Gopher 22"  "Gopher 23"  "Gopher 24"  "Gopher 25" 
 [26] "Gopher 26"  "Gopher 27"  "Gopher 28"  "Gopher 29"  "Gopher 30" 
 [31] "Gopher 31"  "Gopher 32"  "Gopher 33"  "Gopher 34"  "Gopher 35" 
 [36] "Gopher 36"  "Gopher 37"  "Gopher 38"  "Gopher 39"  "Gopher 40" 
 [41] "Gopher 41"  "Gopher 42"  "Gopher 43"  "Gopher 44"  "Gopher 45" 
 [46] "Gopher 46"  "Gopher 47"  "Gopher 48"  "Gopher 49"  "Gopher 50" 
 [51] "Gopher 51"  "Gopher 52"  "Gopher 53"  "Gopher 54"  "Gopher 55" 
 [56] "Gopher 56"  "Gopher 57"  "Gopher 58"  "Gopher 59"  "Gopher 60" 
 [61] "Gopher 61"  "Gopher 62"  "Gopher 63"  "Gopher 64"  "Gopher 65" 
 [66] "Gopher 66"  "Gopher 67"  "Gopher 68"  "Gopher 69"  "Gopher 70" 
 [71] "Gopher 71"  "Gopher 72"  "Gopher 73"  "Gopher 74"  "Gopher 75" 
 [76] "Gopher 76"  "Gopher 77"  "Gopher 78"  "Gopher 79"  "Gopher 80" 
 [81] "Gopher 81"  "Gopher 82"  "Gopher 83"  "Gopher 84"  "Gopher 85" 
 [86] "Gopher 86"  "Gopher 87"  "Gopher 88"  "Gopher 89"  "Gopher 90" 
 [91] "Gopher 91"  "Gopher 92"  "Gopher 93"  "Gopher 94"  "Gopher 95" 
 [96] "Gopher 96"  "Gopher 97"  "Gopher 98"  "Gopher 99"  "Gopher 100"
```
