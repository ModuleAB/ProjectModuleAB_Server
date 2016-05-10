package common

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
	return v.String()
}

func (v version) String() string {
	return fmt.Sprintf("%d.%d.%d-%s", v.main, v.sub, v.build, v.code)
}

var Version version

func init() {
	Version.main = 0
	Version.sub = 0
	Version.build = 0
	Version.code = "Haven"
}
