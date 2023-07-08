package util_error

import (
	"runtime"
	"strings"
	constant "tnals5152/git-pr-slack/const"
)

const DefaultDepth = iota + 2

type CallerInfo struct {
	functionName string
	fileName     string
	line         int
}

func (c *CallerInfo) SetFunctionName(name string) *CallerInfo {
	c.functionName = name
	return c
}

func (c *CallerInfo) SetFileName(fileName string) *CallerInfo {
	c.fileName = fileName
	return c
}

func (c *CallerInfo) SetLine(line int) *CallerInfo {
	c.line = line
	return c
}

type PathError struct {
	err        error
	callerInfo []*CallerInfo
}

func (e *PathError) Error() string {

	if e.err == nil {
		return ""
	}

	return e.err.Error()
}

func (e *PathError) SetError(err error) *PathError {
	e.err = err
	return e
}

func (e *PathError) AddCallerInfo(callerInfo *CallerInfo) *PathError {

	e.callerInfo = append(e.callerInfo, callerInfo)

	return e
}

func wrapError(err error, depth int) {

	pathError, ok := (err).(*PathError)

	if !ok {
		pathError = new(PathError).SetError(err)
	}

	err = pathError

	pc, file, line, ok := runtime.Caller(depth)

	if !ok {
		return
	}

	pcSlice := strings.Split(runtime.FuncForPC(pc).Name(), constant.SLASH)
	fileSlice := strings.Split(file, constant.SLASH)

	pathError.AddCallerInfo(
		new(CallerInfo).
			SetFunctionName(pcSlice[len(pcSlice)-1]).
			SetFileName(fileSlice[len(fileSlice)-1]).
			SetLine(line),
	)

}

// error line 및 에러의 위치를 알고 싶은 함수의 맨 위에 defer DeferWrap(&err)로 작성하면 된다.
func DeferWrap(err error, depths ...int) {
	if err == nil {
		return
	}

	depth := DefaultDepth

	if len(depths) != 0 {
		depth += depths[0]
	}

	wrapError(err, depth)
}
