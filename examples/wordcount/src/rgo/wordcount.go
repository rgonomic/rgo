// Code generated by rgnonomic/rgo; DO NOT EDIT.

package main

/*
#define USE_RINTERNALS
#include <R.h>
#include <Rinternals.h>
extern void R_error(char *s);

// TODO(kortschak): Only emit these when needed.
extern Rboolean Rf_isNull(SEXP s);
extern _GoString_ R_gostring(SEXP x, R_xlen_t i);
extern int getListElementIndex(SEXP list, const char *str);
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/rgonomic/rgo/examples/wordcount"
)

//export Wrapped_Count
func Wrapped_Count(_R_words C.SEXP) C.SEXP {
	defer func() {
		r := recover()
		if r != nil {
			err := C.CString(fmt.Sprint(r))
			C.R_error(err)
			C.free(unsafe.Pointer(err))
		}
	}()

	_p0 := unpackSEXP_types_Slice___string(_R_words)
	_r0 := wordcount.Count(_p0)
	return packSEXP_Count(_r0)
}

func packSEXP_Count(p0 map[string]int) C.SEXP {
	return packSEXP_types_Map_map_string_int(p0)
}

//export Wrapped_CountWithLength
func Wrapped_CountWithLength(_R_words C.SEXP) C.SEXP {
	defer func() {
		r := recover()
		if r != nil {
			err := C.CString(fmt.Sprint(r))
			C.R_error(err)
			C.free(unsafe.Pointer(err))
		}
	}()

	_p0 := unpackSEXP_types_Slice___string(_R_words)
	_r0 := wordcount.CountWithLength(_p0)
	return packSEXP_CountWithLength(_r0)
}

func packSEXP_CountWithLength(p0 map[string]wordcount.WordStats) C.SEXP {
	return packSEXP_types_Map_map_string_github_com_rgonomic_rgo_examples_wordcount_WordStats(p0)
}

//export Wrapped_PrintCountWithLength
func Wrapped_PrintCountWithLength(_R_a C.SEXP) C.SEXP {
	defer func() {
		r := recover()
		if r != nil {
			err := C.CString(fmt.Sprint(r))
			C.R_error(err)
			C.free(unsafe.Pointer(err))
		}
	}()

	_p0 := unpackSEXP_types_Map_map_string_github_com_rgonomic_rgo_examples_wordcount_WordStats(_R_a)
	wordcount.PrintCountWithLength(_p0)
	return C.R_NilValue
}


//export Wrapped_Unique
func Wrapped_Unique(_R_words C.SEXP) C.SEXP {
	defer func() {
		r := recover()
		if r != nil {
			err := C.CString(fmt.Sprint(r))
			C.R_error(err)
			C.free(unsafe.Pointer(err))
		}
	}()

	_p0 := unpackSEXP_types_Slice___string(_R_words)
	_r0 := wordcount.Unique(_p0)
	return packSEXP_Unique(_r0)
}

func packSEXP_Unique(p0 []string) C.SEXP {
	return packSEXP_types_Slice___string(p0)
}

//export Wrapped_UniqueWithLength
func Wrapped_UniqueWithLength(_R_words C.SEXP) C.SEXP {
	defer func() {
		r := recover()
		if r != nil {
			err := C.CString(fmt.Sprint(r))
			C.R_error(err)
			C.free(unsafe.Pointer(err))
		}
	}()

	_p0 := unpackSEXP_types_Slice___string(_R_words)
	_r0 := wordcount.UniqueWithLength(_p0)
	return packSEXP_UniqueWithLength(_r0)
}

func packSEXP_UniqueWithLength(p0 []wordcount.Word) C.SEXP {
	return packSEXP_types_Slice___github_com_rgonomic_rgo_examples_wordcount_Word(p0)
}

//export Wrapped_PrintUniqueWithLength
func Wrapped_PrintUniqueWithLength(_R_a C.SEXP) C.SEXP {
	defer func() {
		r := recover()
		if r != nil {
			err := C.CString(fmt.Sprint(r))
			C.R_error(err)
			C.free(unsafe.Pointer(err))
		}
	}()

	_p0 := unpackSEXP_types_Slice___github_com_rgonomic_rgo_examples_wordcount_Word(_R_a)
	wordcount.PrintUniqueWithLength(_p0)
	return C.R_NilValue
}


func unpackSEXP_types_Basic_int(p C.SEXP) int {
	return int(*C.INTEGER(p))
}

func unpackSEXP_types_Basic_string(p C.SEXP) string {
	return C.R_gostring(p, 0)
}

func unpackSEXP_types_Map_map_string_github_com_rgonomic_rgo_examples_wordcount_WordStats(p C.SEXP) map[string]wordcount.WordStats {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := int(C.Rf_xlength(p))
	r := make(map[string]wordcount.WordStats, n)
	names := C.getAttrib(p, C.R_NamesSymbol)
	for i := 0; i < n; i++ {
		key := string(C.R_gostring(names, C.R_xlen_t(i)))
		r[key] = unpackSEXP_types_Named_github_com_rgonomic_rgo_examples_wordcount_WordStats(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	}
	return r
}

func unpackSEXP_types_Named_github_com_rgonomic_rgo_examples_wordcount_Word(p C.SEXP) wordcount.Word {
	return unpackSEXP_types_Struct_struct_Text_string__rgo___word_____Length_int__rgo___length____(p)
}

func unpackSEXP_types_Named_github_com_rgonomic_rgo_examples_wordcount_WordStats(p C.SEXP) wordcount.WordStats {
	return unpackSEXP_types_Struct_struct_Count_int__rgo___count_____Length_int__rgo___length____(p)
}

func unpackSEXP_types_Slice___github_com_rgonomic_rgo_examples_wordcount_Word(p C.SEXP) []wordcount.Word {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := C.Rf_xlength(p)
	r := make([]wordcount.Word, n)
	for i := range r {
		r[i] = unpackSEXP_types_Named_github_com_rgonomic_rgo_examples_wordcount_Word(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	}
	return r
}

func unpackSEXP_types_Slice___string(p C.SEXP) []string {
	if C.Rf_isNull(p) != 0 {
		return nil
	}
	n := C.Rf_xlength(p)
	r := make([]string, n)
	for i := range r {
		r[i] = string(C.R_gostring(p, C.R_xlen_t(i)))
	}
	return r
}

func unpackSEXP_types_Struct_struct_Count_int__rgo___count_____Length_int__rgo___length____(p C.SEXP) struct{Count int "rgo:\"count\""; Length int "rgo:\"length\""} {
	switch n := C.Rf_xlength(p); {
	case n < 2:
		panic(`missing list element for struct{Count int "rgo:\"count\""; Length int "rgo:\"length\""}`)
	case n > 2:
		err := C.CString(`extra list element ignored for struct{Count int "rgo:\"count\""; Length int "rgo:\"length\""}`)
		C.R_error(err)
		C.free(unsafe.Pointer(err))
	}
	var r struct{Count int "rgo:\"count\""; Length int "rgo:\"length\""}
	var i C.int
	key_count := C.CString("count")
	defer C.free(unsafe.Pointer(key_count))
	i = C.getListElementIndex(p, key_count)
	if i < 0 {
		panic("no list element for field: Count")
	}
	r.Count = unpackSEXP_types_Basic_int(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	key_length := C.CString("length")
	defer C.free(unsafe.Pointer(key_length))
	i = C.getListElementIndex(p, key_length)
	if i < 0 {
		panic("no list element for field: Length")
	}
	r.Length = unpackSEXP_types_Basic_int(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	return r
}

func unpackSEXP_types_Struct_struct_Text_string__rgo___word_____Length_int__rgo___length____(p C.SEXP) struct{Text string "rgo:\"word\""; Length int "rgo:\"length\""} {
	switch n := C.Rf_xlength(p); {
	case n < 2:
		panic(`missing list element for struct{Text string "rgo:\"word\""; Length int "rgo:\"length\""}`)
	case n > 2:
		err := C.CString(`extra list element ignored for struct{Text string "rgo:\"word\""; Length int "rgo:\"length\""}`)
		C.R_error(err)
		C.free(unsafe.Pointer(err))
	}
	var r struct{Text string "rgo:\"word\""; Length int "rgo:\"length\""}
	var i C.int
	key_word := C.CString("word")
	defer C.free(unsafe.Pointer(key_word))
	i = C.getListElementIndex(p, key_word)
	if i < 0 {
		panic("no list element for field: Text")
	}
	r.Text = unpackSEXP_types_Basic_string(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	key_length := C.CString("length")
	defer C.free(unsafe.Pointer(key_length))
	i = C.getListElementIndex(p, key_length)
	if i < 0 {
		panic("no list element for field: Length")
	}
	r.Length = unpackSEXP_types_Basic_int(C.VECTOR_ELT(p, C.R_xlen_t(i)))
	return r
}

func packSEXP_types_Basic_int(p int) C.SEXP {
	return C.ScalarInteger(C.int(p))
}

func packSEXP_types_Basic_string(p string) C.SEXP {
	s := C.Rf_mkCharLenCE(C._GoStringPtr(p), C.int(len(p)), C.CE_UTF8)
	return C.ScalarString(s)
}

func packSEXP_types_Map_map_string_github_com_rgonomic_rgo_examples_wordcount_WordStats(p map[string]wordcount.WordStats) C.SEXP {
	n := len(p)
	r := C.Rf_allocVector(C.VECSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	defer C.Rf_unprotect(1)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	defer C.Rf_unprotect(1)
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		C.SET_VECTOR_ELT(r, i, packSEXP_types_Named_github_com_rgonomic_rgo_examples_wordcount_WordStats(v))
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	return r
}

func packSEXP_types_Map_map_string_int(p map[string]int) C.SEXP {
	n := len(p)
	r := C.Rf_allocVector(C.INTSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	defer C.Rf_unprotect(1)
	names := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(n))
	C.Rf_protect(names)
	defer C.Rf_unprotect(1)
	s := (*[140737488355328]int32)(unsafe.Pointer(C.INTEGER(r)))[:len(p):len(p)]
	var i C.R_xlen_t
	for k, v := range p {
		C.SET_STRING_ELT(names, i, C.Rf_mkCharLenCE(C._GoStringPtr(k), C.int(len(k)), C.CE_UTF8))
		s[i] = int32(v)
		i++
	}
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	return r
}

func packSEXP_types_Named_github_com_rgonomic_rgo_examples_wordcount_Word(p wordcount.Word) C.SEXP {
	return packSEXP_types_Struct_struct_Text_string__rgo___word_____Length_int__rgo___length____(p)
}

func packSEXP_types_Named_github_com_rgonomic_rgo_examples_wordcount_WordStats(p wordcount.WordStats) C.SEXP {
	return packSEXP_types_Struct_struct_Count_int__rgo___count_____Length_int__rgo___length____(p)
}

func packSEXP_types_Slice___github_com_rgonomic_rgo_examples_wordcount_Word(p []wordcount.Word) C.SEXP {
	n := len(p)
	r := C.Rf_allocVector(C.VECSXP, C.R_xlen_t(n))
	C.Rf_protect(r)
	defer C.Rf_unprotect(1)
	for i, v := range p {
		C.SET_VECTOR_ELT(r, C.R_xlen_t(i), packSEXP_types_Named_github_com_rgonomic_rgo_examples_wordcount_Word(v))
	}
	return r
}

func packSEXP_types_Slice___string(p []string) C.SEXP {
	r := C.Rf_allocVector(C.STRSXP, C.R_xlen_t(len(p)))
	C.Rf_protect(r)
	defer C.Rf_unprotect(1)
	for i, v := range p {
		s := C.Rf_mkCharLenCE(C._GoStringPtr(string(v)), C.int(len(v)), C.CE_UTF8)
		C.SET_STRING_ELT(r, C.R_xlen_t(i), s)
	}
	return r
}

func packSEXP_types_Struct_struct_Count_int__rgo___count_____Length_int__rgo___length____(p struct{Count int "rgo:\"count\""; Length int "rgo:\"length\""}) C.SEXP {
	r := C.Rf_allocVector(C.VECSXP, 2)
	C.Rf_protect(r)
	defer C.Rf_unprotect(1)
	names := C.Rf_allocVector(C.STRSXP, 2)
	C.Rf_protect(names)
	defer C.Rf_unprotect(1)
	C.SET_STRING_ELT(names, 0, C.Rf_mkCharLenCE(C._GoStringPtr("count"), 5, C.CE_UTF8))
	C.SET_VECTOR_ELT(r, 0, packSEXP_types_Basic_int(p.Count))
	C.SET_STRING_ELT(names, 1, C.Rf_mkCharLenCE(C._GoStringPtr("length"), 6, C.CE_UTF8))
	C.SET_VECTOR_ELT(r, 1, packSEXP_types_Basic_int(p.Length))
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	return r
}

func packSEXP_types_Struct_struct_Text_string__rgo___word_____Length_int__rgo___length____(p struct{Text string "rgo:\"word\""; Length int "rgo:\"length\""}) C.SEXP {
	r := C.Rf_allocVector(C.VECSXP, 2)
	C.Rf_protect(r)
	defer C.Rf_unprotect(1)
	names := C.Rf_allocVector(C.STRSXP, 2)
	C.Rf_protect(names)
	defer C.Rf_unprotect(1)
	C.SET_STRING_ELT(names, 0, C.Rf_mkCharLenCE(C._GoStringPtr("word"), 4, C.CE_UTF8))
	C.SET_VECTOR_ELT(r, 0, packSEXP_types_Basic_string(p.Text))
	C.SET_STRING_ELT(names, 1, C.Rf_mkCharLenCE(C._GoStringPtr("length"), 6, C.CE_UTF8))
	C.SET_VECTOR_ELT(r, 1, packSEXP_types_Basic_int(p.Length))
	C.setAttrib(r, packSEXP_types_Basic_string("names"), names)
	return r
}

func main() {}
