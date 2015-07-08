// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package build

import (
	"bytes"
	"encoding/hex"
	"fmt"
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

func ContainsString(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetPosLib(value lib, slice []lib) int {
	for i, v := range slice {
		if fmt.Sprintf("%#v", v) == fmt.Sprintf("%#v", value) {
			return i
		}
	}
	return -1
}

func GetHash(i0, i1 int) string {
	return hex.EncodeToString([]byte{byte(i0), byte(i1)})
}
