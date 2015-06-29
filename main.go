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

	dat, err := ioutil.ReadFile("./fog.slib")

	if err != nil {
		fmt.Println("Failed to load file.")
		return
	}

	tokens := make(chan lexer.Token)
	lexer.New("test", string(dat), tokens)

	for token := range tokens {
		fmt.Println(token)
	}
	return

	buildStream := builder.New(tokens)

	for {
		select {
		case err := <-buildStream.Err:
			fmt.Println(err)
			return
		case req := <-buildStream.Request:
			lib, err := ioutil.ReadFile(req.Path)
			if err != nil {
				fmt.Println(fmt.Sprintf("Failed to lib <%s>", req.Path))
				return
			}
			lexer.New(req.Path, string(lib), req.Answer)
		case rep := <-buildStream.Response:
			fmt.Println(rep.GetVertex())
			fmt.Println("-------")
			fmt.Println(rep.GetFragment())
			return
		}
	}
}