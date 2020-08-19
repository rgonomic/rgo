// Copyright ©2019 The rgonomic Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"

	"github.com/rgonomic/rgo/internal/rgo"
	"github.com/rgonomic/rgo/internal/tool"
)

func main() {
	tool.Main(context.Background(), rgo.New(os.Args[0], "", nil), os.Args[1:])
}
