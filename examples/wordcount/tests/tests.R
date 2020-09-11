library(wordcount)

t <- list

failed <- c()
errorf <- function(fmt, ...) {
	msg <- sprintf("%s: %s", deparse(sys.calls()[[sys.nframe()-1]]), sprintf(fmt, ...))
	failed <<- append(failed, msg)
}

printf <- function(fmt, ...) {
	print(sprintf("%s: %s", deparse(sys.calls()[[sys.nframe()-1]]), sprintf(fmt, ...)))
}


count_tests <- list(
	t(name="one two three", input=c("one", "two", "three", "two"), want=c("one"=1L, "three"=1L, "two"=2L)),
NULL)

test_count <- function() {
	for (test in count_tests) {
		if (is.null(test)) {
			next
		}

		got <- wordcount::count(test$input)
		got <- got[order(names(got))] # Map iteration order is unspecified in Go.
		if (any(got != test$want)) {
			errorf("unexpected result for test '%s': got:%s want:%s", test$name, got, test$want)
		}
		if (any(names(got) != names(test$want))) {
			errorf("unexpected names result for test '%s': got:%s want:%s", test$name, names(got), names(test$want))
		}
	}
}
test_count()

test_unique <- function() {
	for (test in count_tests) {
		if (is.null(test)) {
			next
		}

		got <- wordcount::unique(test$input)
		if (any(got != names(test$want))) {
			errorf("unexpected names result for test '%s': got:%s want:%s", test$name, names(got), names(test$want))
		}
	}
}
test_unique()

count_with_length_tests <- list(
	list(name="one two three", input=c("one", "two", "three", "two"), want=list("one"=list("count"=1L, "length"=3L), "three"=list("count"=1L, "length"=5L), "two"=list("count"=2L, "length"=3L))),
NULL)

test_count_with_length <- function() {
	for (test in count_with_length_tests) {
		if (is.null(test)) {
			next
		}

		got <- wordcount::count_with_length(test$input)

		got <- got[order(names(got))] # Map iteration order is unspecified in Go.
		if (any(names(got) != names(test$want))) {
			errorf("unexpected names result for test '%s': got:%s want:%s", test$name, names(got), names(test$want))
			next
		}

		for (n in names(got)) {
			if (any(got[[n]]$count != test$want[[n]]$count)) {
				errorf("unexpected count result for test '%s' element '%s': got:%d want:%d", test$name, n, got[[n]]$count, test$want[[n]]$count)
			}
			if (any(got[[n]]$length != test$want[[n]]$length)) {
				errorf("unexpected length result for test '%s' element '%s': got:%d want:%d", test$name, n, got[[n]]$length, test$want[[n]]$length)
			}
		}

	}
}
test_count_with_length()

test_unique_with_length <- function() {
	for (test in count_with_length_tests) {
		if (is.null(test)) {
			next
		}

		got <- wordcount::count_with_length(test$input)

		got <- got[order(names(got))] # Map iteration order is unspecified in Go.
		if (any(names(got) != names(test$want))) {
			errorf("unexpected names result for test '%s': got:%s want:%s", test$name, names(got), names(test$want))
			next
		}

		for (n in names(got)) {
			if ((got[[n]]$count != test$want[[n]]$count) || (got[[n]]$length != test$want[[n]]$length)) {
				errorf("unexpected result for test '%s' element '%s': got:%s want:%s", test$name, n, got[[n]], test$want[[n]])
			}
		}
	}
}
#test_unique_with_length()

if (length(failed) != 0) {
	print(failed)
	stop("FAILED")
}
