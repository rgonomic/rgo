library(ioutil)

failed <- c()
errorf <- function(fmt, ...) {
	msg <- sprintf("%s: %s", deparse(sys.calls()[[sys.nframe()-1]]), sprintf(fmt, ...))
	failed <<- append(failed, msg)
}

test_read_file <- function() {
	got <- ioutil::read_file("tests.R")
	if (!is.null(got$r1)) {
		errorf("unexpected error: '%s'", got$err)
	}

	if (!startsWith(rawToChar(got$r0), "library(ioutil)")) {
		errorf("unexpected file prefix: got:'%s' want:'library(ioutil)'", rawToChar(got$r0))		
	}
}
test_read_file()

if (length(failed) != 0) {
	print(failed)
	stop("FAILED")
}
