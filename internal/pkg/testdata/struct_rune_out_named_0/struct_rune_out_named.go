// Code generated by "go generate github.com/rgonomic/rgo/internal/pkg/testdata"; DO NOT EDIT.

package struct_rune_out_named_0

//{"out":["int32","struct{F1 int32; F2 int32 \"rgo:\\\"Rname\\\"\"}"]}
func Test0() (res0 struct {
	F1 rune
	F2 rune "rgo:\"Rname\""
}) {
	return res0
}
