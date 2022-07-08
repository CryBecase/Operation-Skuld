package ecode

import (
	"fmt"
	"strings"
)

// E implement Error interface
type E struct {
	code   int
	stderr bool
	msg    string
	args   []interface{}
	tips   tip
}

// New 添加一个不会标准输出的 error
func New(code int, msg string, tip ...string) *E {
	return newError(code, false, msg, tip...)
}

// NewStd 添加一个会标准输出的 error
func NewStd(code int, msg string, tip ...string) *E {
	return newError(code, true, msg, tip...)
}

// newError create a error
func newError(code int, stderr bool, msg string, tip ...string) *E {
	if msg == "" {
		msg = "服务器繁忙"
	}

	tips := ""
	if len(tip) > 0 {
		tips = strings.Join(tip, ",")
	}

	e := &E{
		code:   code,
		stderr: stderr,
		msg:    msg,
		tips:   newTip(tips),
	}

	_em.add(e)

	return e
}

// Error return code in string form
func (e E) Error() string {
	if len(e.args) > 0 {
		return fmt.Sprintf(e.msg, e.args...)
	}
	return e.msg
}

func (e *E) Format(args ...interface{}) *E {
	clone := e.Clone()
	clone.args = args
	return &clone
}

func (e *E) SetMsg(msg string) *E {
	clone := e.Clone()
	clone.msg = msg
	return &clone
}

func (e *E) FormatTips(args ...interface{}) *E {
	clone := e.Clone()
	clone.tips.SetArgs(args)
	return &clone
}

func (e *E) Clone() E {
	return E{
		code:   e.code,
		stderr: e.stderr,
		msg:    e.msg,
		tips: tip{
			tmpl: e.tips.tmpl,
			args: e.tips.args,
		},
	}
}
