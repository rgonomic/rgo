library(cca)

failed <- c()
errorf <- function(fmt, ...) {
	msg <- sprintf("%s: %s", deparse(sys.calls()[[sys.nframe()-1]]), sprintf(fmt, ...))
	failed <<- append(failed, msg)
}

fatalf <- function(fmt, ...) {
	msg <- sprintf("%s: %s", deparse(sys.calls()[[sys.nframe()-1]]), sprintf(fmt, ...))
	stop(msg)
}

mat_list <- function(a) {
	return(list(Rows = nrow(a), Cols = ncol(a), Data = as.vector(a), Stride = nrow(a)))
}

list_mat <- function(a) {
	return(matrix(data = a$Data, nrow = a$Rows, ncol = a$Cols))
}

test_cca <- function() {
	library(MASS)

	x <- cbind(Boston$crim, Boston$indus, Boston$nox, Boston$dis, Boston$rad, Boston$ptratio, Boston$black)
	y <- cbind(Boston$rm, Boston$age, Boston$tax, Boston$medv)
	boston <- cbind(x, y)

	got <- cca::cca(mat_list(x), mat_list(y), NULL)

	if (!is.null(got$err)) {
		fatalf("unexpected error: %s", got$err)
	}

	want_corrs <- c(0.9451239, 1.6786623, 0.5714338, 0.2009740)
	if (norm(as.matrix(got$ccors - want_corrs), "F") > 1e-7) {
		errorf("unexpected correlations: got:%f want:%f", got$corrs, want_corrs)
	}

	want_pVecs <- matrix(nrow=7, ncol=4, data=c(
		-0.25743919,  0.01584775,  0.2122170, -0.09457338,
		-0.48365944,  0.38371019,  0.1474448,  0.65973249,
		-0.08007764,  0.34935567,  0.3287336, -0.28620404,
		 0.12775864, -0.73374277,  0.4851135,  0.22479649,
		-0.69694320, -0.43417488, -0.3602873,  0.02906616,
		-0.09909033,  0.05034112,  0.6384331,  0.10223671,
		 0.42604600,  0.03233344, -0.2289528,  0.64192329
	), byrow=TRUE)
	if (norm(list_mat(got$pVecs) - want_pVecs, "F") > 1e-7) {
		errorf("unexpected pVecs: got:%f want:%f", got$pVecs, want_pVecs)
	}

	want_qVecs <- matrix(nrow=4, ncol=4, data=c(
		 0.01816605, -0.1583489, -0.006672358, -0.98719354,
		-0.23476990,  0.9483315, -0.146242051, -0.15544708,
		-0.97007040, -0.2406072, -0.025183898,  0.02091341,
		 0.05930007, -0.1330460, -0.988905715,  0.02911615
	), byrow=TRUE)
	if (norm(list_mat(got$qVecs) - want_qVecs, "F") > 1e-7) {
		errorf("unexpected qVecs: got:%f want:%f", got$qVecs, want_qVecs)
	}

	want_phiVs <- matrix(nrow=7, ncol=4, data=c(
		-0.0027462234,  0.0093444514,  0.048964393, -0.015496719,
		-0.0428564455, -0.0241708702,  0.036072347,  0.183898323,
		-1.2248435649,  5.6030921365,  5.809414458, -4.792681219,
		-0.0043684825, -0.3424101165,  0.446996122,  0.115016181,
		-0.0741534070, -0.1193135795, -0.111551831,  0.002163876,
		-0.0233270323,  0.1046330818,  0.385304598, -0.016092787,
		 0.0001293051,  0.0004540747, -0.003029632,  0.008189548
	), byrow=TRUE)
	if (norm(list_mat(got$phiVs) - want_phiVs, "F") > 1e-7) {
		errorf("unexpected phiVs: got:%f want:%f", got$phiVs, want_phiVs)
	}

	want_psiVs <- matrix(nrow=4, ncol=4, data=c(
		 0.030159336, -0.30022193,  0.087821738, -1.9583226532,
		-0.006548310,  0.03922121, -0.011757078, -0.0061113064,
		-0.005207552, -0.00457702, -0.002276231,  0.0008441873,
		 0.002011174,  0.00373528, -0.129257807,  0.1037709056
	), byrow=TRUE)
	if (norm(list_mat(got$psiVs) - want_psiVs, "F") > 1e-7) {
		errorf("unexpected psiVs: got:%f want:%f", got$psiVs, want_psiVs)
	}
}
test_cca()

if (length(failed) != 0) {
	print(failed)
	stop("FAILED")
}
