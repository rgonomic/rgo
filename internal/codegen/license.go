// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codegen

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"

	"github.com/rgonomic/rgo/internal/mod"
	"github.com/rgonomic/rgo/internal/pkg"
)

func Licenses(dst FileSystem, prefix string, info *pkg.Info, candidate *regexp.Regexp, verbose bool) error {
	us, err := mod.Module(".")
	if err != nil {
		return fmt.Errorf("failed to get local module information: %w", err)
	}
	them, err := mod.Module(info.Pkg().Path())
	if err != nil {
		return fmt.Errorf("failed to get dependency module information: %w", err)
	}

	licenses, err := mod.Licenses(them.Dir, candidate, verbose)
	if err != nil {
		return fmt.Errorf("failed license scan: %w", err)
	}
	if us == them {
		if verbose && len(licenses) != 0 {
			log.Println("license is local")
		}
		return nil
	}
	for _, l := range licenses {
		path, err := filepath.Rel(them.Dir, l.Path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		if len(licenses) > 1 && prefix != "" {
			path = filepath.Join(prefix, path)
		}
		w, err := dst.Open(path)
		if err != nil {
			return fmt.Errorf("failed open destination file: %w", err)
		}
		_, err = w.Write(l.Text)
		if err != nil {
			return fmt.Errorf("failed to write license file: %w", err)
		}
		w.Close()
	}

	return nil
}
