# R lists and Go maps and structs

This example shows how functions returning maps and slices of non-atomic vectors works with `rgo`.

We have some trivial functions that perform collection of stats on vectors of words.

```
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
```

We do the usual steps to build an `rgo` package, starting with defining the package's `go.mod` file with
```
$ go mod init github.com/rgonomic/rgo/examples/wordcount
```
running `rgo init` with an argument pointing to the current directory since we are at the root of github.com/rgonomic/rgo/examples/cca,
```
$ rgo init .
```
Since there is only one function and it has an all upper-case name, we don't need to make any change to the `rgo.json` file.

The wrapper code is then generated by running the build subcommand.
```
$ rgo build
```
This will generate the Go, C and R wrapper code for the R package. At this stage the `DESCRIPTION` file should be edited.

The package can now be installed.
```
$ R CMD INSTALL .
```

The `count` and `unique` functions return vectors.

```
> library(wordcount)

Attaching package: ‘wordcount’

The following object is masked from ‘package:base’:

    unique

> a <- c("one", "two", "three", "two")
> typeof(wordcount::count(a))
[1] "integer"
> wordcount::count(a)
  one   two three 
    1     2     1 
> typeof(wordcount::unique(a))
[1] "character"
> wordcount::unique(a)
[1] "one"   "three" "two"  
```
The `*_with_length` functions return lists.

```
> typeof(wordcount::count_with_length(a))
[1] "list"
> wordcount::count_with_length(a)
$one
$one$count
[1] 1

$one$length
[1] 3


$two
$two$count
[1] 2

$two$length
[1] 3


$three
$three$count
[1] 1

$three$length
[1] 5


> typeof(wordcount::unique_with_length(a))
[1] "list"
> wordcount::unique_with_length(a)
[[1]]
[[1]]$word
[1] "one"

[[1]]$length
[1] 3


[[2]]
[[2]]$word
[1] "three"

[[2]]$length
[1] 5


[[3]]
[[3]]$word
[1] "two"

[[3]]$length
[1] 3
```

These R objects can be printed in Go syntax, demonstrating the mapping between the types.

```
> print_unique_with_length(unique_with_length(a))
[]wordcount.Word{
 wordcount.Word{
  Text: string("one"),
  Length: int(3),
 },
 wordcount.Word{
  Text: string("three"),
  Length: int(5),
 },
 wordcount.Word{
  Text: string("two"),
  Length: int(3),
 },
}
NULL
> print_count_with_length(count_with_length(a))
map[string]wordcount.WordStats{
 string("one"): wordcount.WordStats{
  Count: int(1),
  Length: int(3),
 },
 string("two"): wordcount.WordStats{
  Count: int(2),
  Length: int(3),
 },
 string("three"): wordcount.WordStats{
  Count: int(1),
  Length: int(5),
 },
}
NULL
```