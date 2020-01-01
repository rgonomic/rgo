// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mod

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/licensecheck"
)

// License holds a packages license information.
type License struct {
	// Path is the path to the license relative
	// to the root of the repository.
	Path string
	// Text is the text of the license.
	Text []byte

	// Cover is the license check information
	// for the license.
	Cover licensecheck.Coverage
}

// Licenses returns all detected licenses under the specified
// root path. If verbose is true, Licenses will log additional
// information to os.Stderr.
func Licenses(root string, candidate *regexp.Regexp, verbose bool) ([]License, error) {
	var files []License
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if candidate != nil && !candidate.MatchString(info.Name()) {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		cover, ok := licensecheck.Cover(b, licensecheck.Options{})
		if ok {
			files = append(files, License{Path: path, Cover: cover, Text: b})
		} else if verbose {
			log.Println("missed:", path)
		}
		return nil
	})
	return files, err
}
