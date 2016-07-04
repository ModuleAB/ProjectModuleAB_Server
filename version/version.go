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

var Version version

func init() {
	Version = version{
		main:  0,
		sub:   0,
		build: 167,
		code:  "Husky",
	}
}
