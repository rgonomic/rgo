#!/bin/bash

status=0

cd examples
R CMD build * || status=1
R CMD check --no-manual *.tar.gz || status=1
rm -rf *.Rcheck *.tar.gz

exit $status
