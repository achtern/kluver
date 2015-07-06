// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/achtern/kluver/compiler"
)

func main() {

	includePath := flag.String("include-path", "./", "path prefix for all imported libs")
	flag.Parse()

	tail := flag.Args();
	
	if len(tail) == 0 {
		fmt.Fprintf(os.Stderr, "missing path to shader source file\n")
		os.Exit(1)
	}


	shader, err := compiler.New("example.shader", *includePath, compiler.FileLoader{})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("VERTEX--")
	fmt.Println(shader.GetVertex())
	fmt.Println("FRAGMENT--")
	fmt.Println(shader.GetFragment())
}
