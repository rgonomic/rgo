// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mod

import (
	"bytes"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

// The Go LICENSE is the only license that we can guarantee will be on the system.
const goLicense = `Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

func TestLicenses(t *testing.T) {
	candidate := regexp.MustCompile("LICENSE")
	licenses, err := Licenses(runtime.GOROOT(), candidate, false)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	var foundLicense bool
	var foundPath bool
	for _, l := range licenses {
		if bytes.Equal(l.Text, []byte(goLicense)) {
			foundLicense = true
			path, err := filepath.Rel(runtime.GOROOT(), l.Path)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if path == "LICENSE" {
				foundPath = true
			}
		}
	}
	if !foundLicense {
		t.Error("failed to find Go LICENSE text in GOROOT")
	}
	if !foundPath {
		t.Error("failed to find Go LICENSE path in GOROOT")
	}
}
