// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package util

import (
	"strconv"
	"strings"
)

func GetLineFromPos(input, pos string) int {
	upper, _ := strconv.Atoi(pos)
	// add one, since line numbers do not start at zero
	return strings.Count(input[0:upper], "\n") + 1
}

func AddTrailingSlash(input string) string {
	if strings.HasSuffix(input, "/") {
		return input
	}
	return input + "/"
}

func GetSubStringUpTo(input string, breaker rune) string {
	for i, r := range input {
		if r == breaker {
			return input[:i]
		}
	}

	return input
}

func ConstructFileName(path, name, suffix string) string {
	path = AddTrailingSlash(path)
	return path + name + suffix
}