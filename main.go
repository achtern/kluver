// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"

	builder "github.com/achtern/kluver/build"
	"github.com/achtern/kluver/lexer"
)

func main() {

	dat, err := ioutil.ReadFile("./example.shader")

	if err != nil {
		fmt.Println("Failed to load file.")
		return
	}

	tokens := make(chan lexer.Token)
	lexer.New("test", string(dat), tokens)

	buildStream := builder.New(tokens)

	for {
		select {
		case err := <-buildStream.Err:
			panic(err)
		case req := <-buildStream.Request:
			lib, err := ioutil.ReadFile(req.Path)
			if err != nil {
				fmt.Println("Failed to load lib file.")
				return
			}
			lexer.New(req.Path, string(lib), req.Answer)
		case rep := <-buildStream.Response:
			fmt.Println(rep.String())
			fmt.Println("-------")
			fmt.Println(rep.String())
		}
	}
}
