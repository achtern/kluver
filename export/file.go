// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package export

import (
	"io/ioutil"
)

func WriteFile(contents, path string) error {
	return ioutil.WriteFile(path, []byte(contents), 0644)
}
