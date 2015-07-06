// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"errors"
	builder "github.com/achtern/kluver/build"
	"github.com/achtern/kluver/lexer"
	"github.com/achtern/kluver/util"
)

func New(path, includePath string,loader FileLoader) (*builder.Shader, error) {
	raw, err := loader.Get(path)
	if err != nil {
		return nil, err
	}

	// prepare includePath
	includePath = util.AddTrailingSlash(includePath)

	// initial
	tokens := make(chan lexer.Token)
	lexer.New(path, raw, tokens)

	buildStream := builder.New(tokens)

	for {
		select {
		case err := <-buildStream.Err:
			return nil, err
		case req := <-buildStream.Request:
			lib, err := loader.Get(includePath + req.Path)
			if err != nil {
				return nil, errors.New("Failed to lib <" + req.Path + ">")
			}
			lexer.New(req.Path, lib, req.Answer)
		case rep := <-buildStream.Response:
			return &rep, nil
		}
	}

}