// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package camel implements camel case word splitting.
package camel

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Split returns a slice of strings split according to camel case rules.
// The string is split so that words start with a leading capital unless
// it is the first word in the string, which may start with a lower case.
// The underscore character breaks words and numbers are separated from
// words.
// Invalid UTF-8 encoded strings are not split.
func Split(str string) []string {
	if str == "" || !utf8.ValidString(str) {
		return []string{str}
	}
	var words []string
	var i, last int
	var prev, curr rune
	for i, curr = range str {
		switch {
		case curr == '_':
			if prev != curr && last != i {
				words = append(words, str[last:i])
			}
			last = i + 1

		case unicode.IsNumber(curr):
			if !unicode.IsNumber(prev) && last != i {
				words = append(words, str[last:i])
				last = i
			}

		case unicode.IsUpper(curr):
			next, _ := utf8.DecodeRuneInString(str[i+utf8.RuneLen(curr):])
			if unicode.IsLower(prev) || unicode.IsNumber(prev) || (unicode.IsUpper(prev) && unicode.IsLower(next)) {
				words = append(words, str[last:i])
				last = i
			}
		}

		prev = curr
	}
	if last < len(str) {
		words = append(words, str[last:])
	}
	return words
}

// Splitter implements camel case word splitting with additional control
// over word boundaries where case-sensitive known words can be provided.
// The zero value of Splitter can be used and Splitter{}.Split is the
// equivalent of Split.
type Splitter struct {
	known   []string
	isKnown map[string]bool
}

// NewSplitter returns a new camel case splitter using the provided known
// words. The words should be provided in order of checking priority.
func NewSplitter(known []string) Splitter {
	s := Splitter{
		known:   known,
		isKnown: make(map[string]bool),
	}
	for _, w := range known {
		s.isKnown[w] = true
	}
	return s
}

// Split returns a slice of strings split according to camel case rules
// after having first split the string according the provided list of known
// words in the order they were provided to NewSplitter. When known words
// overlap in the string, only the first matching word is used.
func (s Splitter) Split(str string) []string {
	if str == "" || !utf8.ValidString(str) {
		return []string{str}
	}
	if len(s.known) == 0 {
		return Split(str)
	}
	var words []string
	for _, w := range splitKnownWords(str, s.known) {
		if s.isKnown[w] {
			words = append(words, w)
		} else {
			words = append(words, Split(w)...)
		}
	}
	return words
}

func splitKnownWords(str string, words []string) []string {
	if str == "" {
		return nil
	}
	if len(words) == 0 {
		return []string{str}
	}
	w := words[0]
	i := strings.Index(str, w)
	if i < 0 {
		return splitKnownWords(str, words[1:])
	}
	front := append(splitKnownWords(str[0:i], words[1:]),
		str[i:i+len(w)])
	back := splitKnownWords(str[i+len(w):], words)
	if len(back) == 0 {
		return front
	}
	return append(front, back...)
}
