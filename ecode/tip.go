package ecode

import (
	"fmt"
)

type tip struct {
	tmpl string
	args []interface{}
}

func newTip(s string) tip {
	return tip{
		tmpl: s,
	}
}

func (t *tip) String() string {
	if len(t.args) > 0 {
		return fmt.Sprintf(t.tmpl, t.args...)
	}
	return t.tmpl
}

func (t *tip) SetArgs(a ...interface{}) {
	if len(a) == 0 {
		return
	}
	t.args = a
}
