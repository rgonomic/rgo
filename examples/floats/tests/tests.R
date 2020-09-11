library(floats)

failed <- c()
errorf <- function(fmt, ...) {
	msg <- sprintf("%s: %s", deparse(sys.calls()[[sys.nframe()-1]]), sprintf(fmt, ...))
	failed <<- append(failed, msg)
}

test_cum_prod <- function() {
	set.seed(1)
	for (i in c(1L:10L)) {
		input <- rnorm(10)
		want <- cumprod(input)
		dst <- vector(mode="double", length=10)
		got <- floats::cum_prod(dst, input)
		if (norm(as.matrix(got - want), "F") > 1e-14) {
			errorf("unexpected cumulative product for tes %d: got:%f want:%f", i, got, want)
		}
	}
}
test_cum_prod()

test_cum_sum <- function() {
	set.seed(1)
	for (i in c(1L:10L)) {
		input <- rnorm(10)
		want <- cumsum(input)
		dst <- vector(mode="double", length=10)
		got <- floats::cum_sum(dst, input)
		if (norm(as.matrix(got - want), "F") > 1e-14) {
			errorf("unexpected cumulative sum for tes %d: got:%f want:%f", i, got, want)
		}
	}
}
test_cum_sum()

if (length(failed) != 0) {
	print(failed)
	stop("FAILED")
}
