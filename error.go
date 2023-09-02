package errcode

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// Error 错误
type Error struct {
	code  int
	msg   string
	err   error
	stack *stack
	data  interface{}
}

// stack 堆栈
type stack struct {
	file string
	fn   string
	line int
}

// New 实例化 Error
func New(code int, msg string) *Error {
	return newError(code, msg)
}

func newError(code int, msg string) *Error {
	e := &Error{
		code: code,
		msg:  msg,
	}
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return e
	}

	fn := runtime.FuncForPC(pc)
	funcName := fn.Name()
	s := strings.Split(funcName, ".")
	if len(s) > 0 {
		funcName = s[len(s)-1]
	}

	e.stack = &stack{
		file: file,
		fn:   funcName,
		line: line,
	}

	return e
}

// Wrap 包装错误
func (e *Error) Wrap(err error) *Error {
	e.err = err

	return e
}

// WithData 添加数据
func (e *Error) WithData(data interface{}) *Error {
	e.data = data
	return e
}

// Code 获取错误码
func (e *Error) Code() int {
	return e.code
}

// Msg 获取错误消息
func (e *Error) Msg() string {
	return e.msg
}

// Data 获取数据
func (e *Error) Data() interface{} {
	return e.data
}

// Unwrap 解包错误
// 实现 interface Unwrap
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

// AddMsg 追加 msg
func (e *Error) AddMsg(msg string, v ...interface{}) *Error {
	if len(msg) > 0 {
		e.msg += msg
	}

	if len(v) > 0 {
		e.msg = fmt.Sprintf(" "+e.msg, v...)
	}

	return e
}

// SetMsg 设置 msg
func (e *Error) SetMsg(msg string, v ...interface{}) *Error {
	e.msg = msg

	if len(v) > 0 {
		e.msg = fmt.Sprintf(e.msg, v...)
	}

	return e
}

// FillMsg 填充 msg
func (e *Error) FillMsg(v ...interface{}) *Error {
	if len(v) > 0 {
		e.msg = fmt.Sprintf(e.msg, v...)
	}

	return e
}

// Is 判断错误是否相对（错误码相同则认为相等）
func (e *Error) Is(target error) bool {
	if e == nil {
		return e == target
	}

	t, ok := target.(*Error)
	if !ok {
		return false
	}
	if t == nil {
		return e == t
	}

	if e.code != t.code {
		return false
	}

	return true
}

// Error 错误字符串
// 实现 interface error
func (e *Error) Error() string {
	var s string
	if e.stack != nil {
		s = "file=" + e.stack.file + ":" + e.stack.fn + ":" + strconv.Itoa(e.stack.line) + ", "
	}

	s += "code=" + strconv.Itoa(e.code) + ", msg=" + e.msg

	if e.err != nil {
		s += ", cause=" + e.err.Error()
	}

	if e.data != nil {
		s = fmt.Sprintf("%s, data=%v", s, e.data)
	}

	return s
}
