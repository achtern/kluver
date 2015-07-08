// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package build

import (
	"bytes"
	"encoding/hex"
	_ "github.com/achtern/kluver/lexer"
	"strings"
)

type StringBuffer struct {
	buffer bytes.Buffer
}

func (sb *StringBuffer) Append(s string) {
	sb.buffer.WriteString(s)
}

func (sb *StringBuffer) String() string {
	return sb.buffer.String()
}

func Contains(s []Tokens, e Tokens) bool {
	for _, tokens := range s {
		typ := strings.Trim(tokens[0].Val, " ") == strings.Trim(e[0].Val, " ")
		val := strings.Trim(tokens[1].Val, " ") == strings.Trim(e[1].Val, " ")
		if typ && val {
			return true
		}
	}
	return false
}

func GetHash(i0, i1 int) string {
	return hex.EncodeToString([]byte{byte(i0), byte(i1)})
}