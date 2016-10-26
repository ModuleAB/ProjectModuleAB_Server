/*ModuleAB version/version.go -- define version
 * Copyright (C) 2016 TonyChyi <tonychee1989@gmail.com>
 * License: GPL v3 or later.
 */

package version

import (
	"fmt"
)

type version struct {
	main  int
	sub   int
	build int
	code  string
}

func (v version) GetVersion() string {
	return fmt.Sprintf("%d.%d.%d", v.main, v.sub, v.build)
}

func (v version) GetCode() string {
	return v.code
}

func (v version) String() string {
	return fmt.Sprintf("%s-%s", v.GetVersion(), v.GetCode())
}

//Version is global version instance.
var Version version

func init() {
	Version = version{
		main:  0,
		sub:   0,
		build: 154,
		code:  "Husky",
	}
}
