package errcode

import (
	"fmt"
)

// Code 状态码
type Code struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Code 状态码转换为错误
func (c Code) Error(v ...interface{}) *Error {
	if len(v) > 0 {
		c.Msg = fmt.Sprintf(c.Msg, v...)
	}

	return newError(c.Code, c.Msg)
}
