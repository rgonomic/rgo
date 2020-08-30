package wordcount

import (
	"sort"

	"github.com/kortschak/utter"
)

// Count returns a list of counts of occurrences of words.
func Count(words []string) map[string]int {
	m := make(map[string]int)
	for _, w := range words {
		m[w]++
	}
	return m
}

type WordStats struct {
	Count  int `rgo:"count"`
	Length int `rgo:"length"`
}

// CountWithLength returns a list of counts of occurrences of words,
// redundantly keeping the word length.
func CountWithLength(words []string) map[string]WordStats {
	m := make(map[string]WordStats)
	for _, w := range words {
		s := m[w]
		s.Count = len(w)
		s.Length++
		m[w] = s
	}
	return m
}

// PrintUniqueWithLength prints the Go syntax representation of a.
func PrintCountWithLength(a map[string]WordStats) { utter.Dump(a) }

// Unique returns a vector of unique words from the words input vector
// sorted ascending.
func Unique(words []string) []string {
	m := Count(words)
	s := make([]string, 0, len(m))
	for w := range m {
		s = append(s, w)
	}
	sort.Strings(s)
	return s
}

type Word struct {
	Text   string `rgo:"word"`
	Length int    `rgo:"length"`
}

// UniqueWithLength returns a vector of unique words from the words input
// vector sorted ascending, redundantly keeping the word length.
func UniqueWithLength(words []string) []Word {
	m := Count(words)
	s := make([]Word, 0, len(m))
	for w := range m {
		s = append(s, Word{Text: w, Length: len(w)})
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].Text < s[j].Text
	})
	return s
}

// PrintUniqueWithLength prints the Go syntax representation of a.
func PrintUniqueWithLength(a []Word) { utter.Dump(a) }
