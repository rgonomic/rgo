package sexp

import (
	"fmt"
	"unsafe"

	"github.com/rgonomic/rgo/sexp"
)

// PrintSEXP prints the SEXP value passed to it along with any names it has
// and returns it unaltered. If the value is an atomic vector, its contents
// are printed.
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
