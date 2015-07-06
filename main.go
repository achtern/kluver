// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/achtern/kluver/compiler"
	"github.com/achtern/kluver/util"
	"github.com/achtern/kluver/export"
)

func main() {

	includePath := flag.String("include-path", "./", "path prefix for all imported libs")
	vertexTargetSuffix := flag.String("vertex-target-suffix", "_vert.glsl", "suffix for the generated vertex shader")
	fragmentTargetSuffix := flag.String("fragment-target-suffix", "_frag.glsl", "suffix for the generated fragment shader")
	exportPath := flag.String("export-path", "./", "destination path for generated glsl files")
	flag.Parse()

	tail := flag.Args();
	
	if len(tail) == 0 {
		fmt.Fprintf(os.Stderr, "missing path to shader source file\n")
		os.Exit(1)
	}

	shaderSourcePath := tail[0]


	shader, err := compiler.New(shaderSourcePath, *includePath, compiler.FileLoader{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%q\n", err)
		os.Exit(1)
	}

	// Prepare exportPath
	cleanExportPath := util.AddTrailingSlash(*exportPath)

	vertErr := export.WriteFile(shader.GetVertex(), cleanExportPath + "example" + *vertexTargetSuffix)
	fragErr := export.WriteFile(shader.GetFragment(), cleanExportPath + "example" + *fragmentTargetSuffix)

	if vertErr != nil {
		fmt.Fprintf(os.Stderr, "error writing vertex destination file:\n%q\n", vertErr)
	}
	if fragErr != nil {
		fmt.Fprintf(os.Stderr, "error writing fragment destination file:\n%q\n", fragErr)
	}
	if vertErr != nil || fragErr != nil {
		os.Exit(1)
	}
}
