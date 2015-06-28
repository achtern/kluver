// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package build

import (
	"bytes"
)

type StringBuffer struct {
	buffer bytes.Buffer
}

func (sb *StringBuffer) append(s string) {
	sb.buffer.WriteString(s)
}

func (sb *StringBuffer) String() string {
	return sb.buffer.String()
}
