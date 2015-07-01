// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"io/ioutil"
)

type FileLoader struct {
}

func (f *FileLoader) Get(path string) (output string, err error) {
	dat, err := ioutil.ReadFile(path)
	return string(dat), err
}