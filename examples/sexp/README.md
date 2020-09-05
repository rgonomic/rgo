# Direct access to SEXP values

Brave users can directly access the R `SEXP` value by accepting/returning `unsafe.Pointer` parameters (see the [sexp package docs](https://pkg.go.dev/github.com/rgonomic/rgo/sexp?tab=doc) for more details).

```
// PrintSEXP prints the SEXP value passed to it and returns it unaltered.
// If the value is an atomic vector, its contents are printed.
func PrintSEXP(p unsafe.Pointer) unsafe.Pointer {
	sxp := (*sexp.Value)(p)
	info := sxp.Info()
	switch sxp := sxp.Valuer().(type) {
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
// Gophers returns n gophers with name attributes.
func Gophers(n int) unsafe.Pointer {
	c := sexp.NewString(n).Protect()
	defer c.Unprotect()
	names := sexp.NewString(n).Protect()
	defer names.Unprotect()

	vec := c.Vector()
	namesVec := names.Vector()
	for i := range vec {
		vec[i] = sexp.NewCharacter(fmt.Sprintf("Gopher %d", i+1))
		namesVec[i] = sexp.NewCharacter(fmt.Sprintf("Name_%d", i+1))
	}

	sxp := c.Value()
	sxp.SetNames(names)
	return sxp.Export()
}
```

```
> gophers(1L)
    Name_1 
"Gopher 1" 
> gophers(100L)
      Name_1       Name_2       Name_3       Name_4       Name_5       Name_6 
  "Gopher 1"   "Gopher 2"   "Gopher 3"   "Gopher 4"   "Gopher 5"   "Gopher 6" 
      Name_7       Name_8       Name_9      Name_10      Name_11      Name_12 
  "Gopher 7"   "Gopher 8"   "Gopher 9"  "Gopher 10"  "Gopher 11"  "Gopher 12" 
     Name_13      Name_14      Name_15      Name_16      Name_17      Name_18 
 "Gopher 13"  "Gopher 14"  "Gopher 15"  "Gopher 16"  "Gopher 17"  "Gopher 18" 
     Name_19      Name_20      Name_21      Name_22      Name_23      Name_24 
 "Gopher 19"  "Gopher 20"  "Gopher 21"  "Gopher 22"  "Gopher 23"  "Gopher 24" 
     Name_25      Name_26      Name_27      Name_28      Name_29      Name_30 
 "Gopher 25"  "Gopher 26"  "Gopher 27"  "Gopher 28"  "Gopher 29"  "Gopher 30" 
     Name_31      Name_32      Name_33      Name_34      Name_35      Name_36 
 "Gopher 31"  "Gopher 32"  "Gopher 33"  "Gopher 34"  "Gopher 35"  "Gopher 36" 
     Name_37      Name_38      Name_39      Name_40      Name_41      Name_42 
 "Gopher 37"  "Gopher 38"  "Gopher 39"  "Gopher 40"  "Gopher 41"  "Gopher 42" 
     Name_43      Name_44      Name_45      Name_46      Name_47      Name_48 
 "Gopher 43"  "Gopher 44"  "Gopher 45"  "Gopher 46"  "Gopher 47"  "Gopher 48" 
     Name_49      Name_50      Name_51      Name_52      Name_53      Name_54 
 "Gopher 49"  "Gopher 50"  "Gopher 51"  "Gopher 52"  "Gopher 53"  "Gopher 54" 
     Name_55      Name_56      Name_57      Name_58      Name_59      Name_60 
 "Gopher 55"  "Gopher 56"  "Gopher 57"  "Gopher 58"  "Gopher 59"  "Gopher 60" 
     Name_61      Name_62      Name_63      Name_64      Name_65      Name_66 
 "Gopher 61"  "Gopher 62"  "Gopher 63"  "Gopher 64"  "Gopher 65"  "Gopher 66" 
     Name_67      Name_68      Name_69      Name_70      Name_71      Name_72 
 "Gopher 67"  "Gopher 68"  "Gopher 69"  "Gopher 70"  "Gopher 71"  "Gopher 72" 
     Name_73      Name_74      Name_75      Name_76      Name_77      Name_78 
 "Gopher 73"  "Gopher 74"  "Gopher 75"  "Gopher 76"  "Gopher 77"  "Gopher 78" 
     Name_79      Name_80      Name_81      Name_82      Name_83      Name_84 
 "Gopher 79"  "Gopher 80"  "Gopher 81"  "Gopher 82"  "Gopher 83"  "Gopher 84" 
     Name_85      Name_86      Name_87      Name_88      Name_89      Name_90 
 "Gopher 85"  "Gopher 86"  "Gopher 87"  "Gopher 88"  "Gopher 89"  "Gopher 90" 
     Name_91      Name_92      Name_93      Name_94      Name_95      Name_96 
 "Gopher 91"  "Gopher 92"  "Gopher 93"  "Gopher 94"  "Gopher 95"  "Gopher 96" 
     Name_97      Name_98      Name_99     Name_100 
 "Gopher 97"  "Gopher 98"  "Gopher 99" "Gopher 100" 
```
