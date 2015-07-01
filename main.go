// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/achtern/kluver/compiler"
)

func main() {
	shader, err := compiler.New("example.shader", compiler.FileLoader{})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("VERTEX--")
	fmt.Println(shader.GetVertex())
	fmt.Println("FRAGMENT--")
	fmt.Println(shader.GetFragment())
}
