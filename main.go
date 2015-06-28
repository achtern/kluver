// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/achtern/kluver/lexer"
	builder "github.com/achtern/kluver/build"
	"io/ioutil"
)

func main() {

	dat, err := ioutil.ReadFile("./example.shader")

	if err != nil {
		fmt.Println("Failed to load file.")
	}

	_, tokens := lexer.Lex("test", string(dat))

	vertex, fragment, _ := builder.Build(tokens)
	fmt.Println(vertex)
	fmt.Println("-------")
	fmt.Println(fragment)
}
