// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/achtern/kluver/compiler"
	"github.com/achtern/kluver/export"
)

func main() {

	includePath := flag.String("include-path", "./", "path prefix for all imported libs")
	vertexTargetSuffix := flag.String("vertex-target-suffix", ".gvs", "suffix for the generated vertex shader")
	fragmentTargetSuffix := flag.String("fragment-target-suffix", ".gfs", "suffix for the generated fragment shader")
	flag.Parse()

	tail := flag.Args();
	
	if len(tail) == 0 {
		fmt.Fprintf(os.Stderr, "missing path to shader source file\n")
		os.Exit(1)
	}

	shaderSourcePath := tail[0]


	shader, err := compiler.New(shaderSourcePath, *includePath, compiler.FileLoader{})
	if err != nil {
		fmt.Println(err)
		return
	}

	export.WriteFile(shader.GetVertex(), "example" + *vertexTargetSuffix)
	export.WriteFile(shader.GetFragment(), "example" + *fragmentTargetSuffix)
}
