library(sexp)

t <- list

failed <- c()
errorf <- function(fmt, ...) {
	msg <- sprintf("%s: %s", deparse(sys.calls()[[sys.nframe()-1]]), sprintf(fmt, ...))
	failed <<- append(failed, msg)
}

fatalf <- function(fmt, ...) {
	msg <- sprintf("%s: %s", deparse(sys.calls()[[sys.nframe()-1]]), sprintf(fmt, ...))
	stop(msg)
}

print_sexp_tests <- list(
	t(name="floating point scalar", input=1, want=1),
	t(name="floating point vector", input=c(1, 2), want=c(1, 2)),
	t(name="integer scalar", input=1L, want=1L),
	t(name="integer vector", input=c(1L, 2L), want=c(1L, 2L)),
	t(name="string scalar", input="a", want="a"),
	t(name="string vector", input=c("a", "b"), want=c("a", "b")),
	t(name="complex scalar", input=0+1i, want=0+1i),
	t(name="complex vector", input=c(0+1i, 1+1i), want=c(0+1i, 1+1i)),
NULL)

test_print_sexp <- function() {
	for (test in print_sexp_tests) {
		if (is.null(test)) {
			next
		}

		got <- sexp::print_sexp(test$input)
		if (got != test$want) {
			errorf("unexpected result for test '%s': got:%s want:%s", test$name, got, test$want)
		}
	}
}
test_print_sexp()

gophers_tests <- list(
	t(name="One gopher", input=1L, want=c("Gopher 1"), want_names=c("Name_1")),
	t(name="Two gophers", input=2L, want=c("Gopher 1", "Gopher 2"), want_names=c("Name_1", "Name_2")),
NULL)

test_gophers <- function() {
	for (test in gophers_tests) {
		if (is.null(test)) {
			next
		}

		got <- sexp::gophers(test$input)
		if (got != test$want) {
			errorf("unexpected result for test '%s': got:%s want:%s", test$name, got, test$want)
		}
		if (names(got) != test$want_names) {
			errorf("unexpected names for test '%s': got:%s want:%s", test$name, names(got), test$want_names)
		}
	}
}
test_gophers()

if (length(failed) != 0) {
	print(failed)
	stop("FAILED")
}
