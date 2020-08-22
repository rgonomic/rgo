// Copyright ©2020 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codegen

import (
	"fmt"

	"golang.org/x/mod/semver"

	"github.com/rgonomic/rgo/internal/mod"
	"github.com/rgonomic/rgo/internal/pkg"
)

func Description(dst FileSystem, info *pkg.Info) error {
	origin, err := mod.Module(info.Pkg().Path())
	if err != nil {
		return fmt.Errorf("failed to get dependency module information: %w", err)
	}

	w, err := dst.Open("DESCRIPTION")
	if err != nil {
		return fmt.Errorf("failed open destination file: %w", err)
	}
	_, err = fmt.Fprintf(w, `Package: %s
Title: What the Package Does (One Line, Title Case)
Version: %s
Authors@R:
    person(given   = "First",
           family  = "Last",
           role    = c("aut", "cre"),
           email   = "first.last@example.com",
           comment = c(ORCID = "YOUR-ORCID-ID"))
Description: What the package does (one paragraph).
License: See LICENSE directory
Encoding: UTF-8
LazyData: true
`, info.Pkg().Name(), version(origin.Version))
	if err != nil {
		return fmt.Errorf("failed to write DESCRIPTION file: %w", err)
	}
	w.Close()

	return nil
}

func version(v string) string {
	if v == "" {
		return "0.0.0"
	}
	if semver.IsValid(v) {
		return v[1:]
	}
	return v
}
