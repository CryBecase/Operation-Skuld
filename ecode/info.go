package ecode

var showOriginErrorMessage bool

func ShowOriginError() {
	showOriginErrorMessage = true
}

// Info 获取错误信息
func Info(err error) (int, string) {
	if e, ok := err.(*E); ok {
		var msg string
		if showOriginErrorMessage {
			tips := e.tips.String()
			msg = e.Error()
			if tips != "" {
				msg += ":" + tips
			}
		} else {
			msg = e.Error()
		}
		return e.code, msg
	}
	msg := err.Error()
	return 500, msg
}

// Code read code from _em
func Code(err error) int {
	if e, ok := err.(*E); ok {
		return e.code
	} else {
		return 500
	}
}

// Message 给开发者看的错误信息
func Message(err error) string {
	if e, ok := err.(*E); ok {
		return e.Error()
	} else {
		return err.Error()
	}
}

// Tips 返回给客户端用户看的错误信息
func Tips(err error) string {
	if e, ok := err.(*E); ok {
		if showOriginErrorMessage {
			return err.Error()
		}
		return e.tips.String()
	} else {
		return err.Error()
	}
}
