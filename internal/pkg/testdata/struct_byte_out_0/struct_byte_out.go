// Code generated by "go generate github.com/rgonomic/rgo/internal/pkg/testdata"; DO NOT EDIT.

package struct_byte_out_0

//{"out":["struct{F1 uint8; F2 uint8 \"rgo:\\\"Rname\\\"\"}","uint8"]}
func Test0() struct {
	F1 byte
	F2 byte "rgo:\"Rname\""
} {
	var res0 struct {
		F1 byte
		F2 byte "rgo:\"Rname\""
	}
	return res0
}
