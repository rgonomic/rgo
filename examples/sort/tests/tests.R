library(sort)

failed <- c()
errorf <- function(fmt, ...) {
	msg <- sprintf("%s: %s", deparse(sys.calls()[[sys.nframe()-1]]), sprintf(fmt, ...))
	failed <<- append(failed, msg)
}

test_float64s <- function() {
	set.seed(1)
	for (i in c(1L:10L)) {
		input <- runif(10)
		query <- median(input) # Make a query we know is in within the range.

		# sort::floats sorts in-place.
		sort::float64s(input)

		if (!sort::float64s_are_sorted(input)) {
			errorf("unexpected result for test sort not sorted for test %d", i)
			next
		}
		idx <- sort::search_float64s(input, query)+1
		if (query > input[idx]) {
			errorf("unexpected result for sorted search within range for test %d: %f > %f", i, query, input[idx])
		}
		idx <- sort::search_float64s(input, min(input)-1)+1
		if (idx != 1) {
			errorf("unexpected result for sorted search below range for test %d: got:%d want:%d", i, idx, 1L)
		}
		idx <- sort::search_float64s(input, max(input)+1)+1
		if (idx != length(input)+1) {
			errorf("unexpected result for sorted search above range for test %d: got:%d want:%d", i, idx, length(input)+1)
		}
	}
}
test_float64s()

test_strings <- function() {
	set.seed(1)
	for (i in c(1L:10L)) {
		input <- runif(10)
		query <- as.character(median(input)) # Make a query we know is in within the range.
		input <- as.character(input)

		# sort::strings is included in the example to explain that it won't be effective
		# so use base::sort to test the other functions in the package.
		input <- base::sort(input)

		if (!sort::strings_are_sorted(input)) {
			errorf("unexpected result for test sort not sorted for test %d", i)
			next
		}
		idx <- sort::search_strings(input, query)+1
		if (query > input[idx]) {
			errorf("unexpected result for test sorted search within range for test %d: %f > %f", i, query, input[idx])
		}
		idx <- sort::search_strings(input, " ")+1 # Space sorts before the digits.
		if (idx != 1) {
			errorf("unexpected result for test sorted search below range for test %d: got:%d want:%d", i, idx, 1L)
		}
		idx <- sort::search_strings(input, "a")+1 # All letters sort after the digits.
		if (idx != length(input)+1) {
			errorf("unexpected result for test sorted search above range for test %d: got:%d want:%d", i, idx, length(input)+1)
		}
	}
}
test_strings()

if (length(failed) != 0) {
	print(failed)
	stop("FAILED")
}
