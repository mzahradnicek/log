package log

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

const (
	TError         = "error"
	TWarning       = "warning"
	TInfo          = "info"
	TDebug         = "debug"
	maxStackLength = 50
)

type Entrier interface {
	ToFields() Fields
	Parent() error
	Unwrap() error
}

type Fields map[string]interface{}

type logEntry struct {
	ltype      string
	message    string
	line       int
	file       string
	fields     Fields
	stackTrace []string

	err error
}

func (e *logEntry) Error() string {
	if e.ltype == "error" {
		return e.message
	}

	return ""
}

func (e *logEntry) Unwrap() error {
	return e.err
}

// Deprecated
func (e *logEntry) Parent() error {
	return e.err
}

// Deprecated
func GetParentError(e error) error {
	if ent, ok := e.(Entrier); ok {
		return ent.Parent()
	}
	return e
}

// ToMap convert logEntry to map for use in Log function.
func (e *logEntry) ToFields() Fields {
	res := Fields{}
	for k, v := range e.fields {
		res[k] = v
	}

	res["type"] = e.ltype
	res["msg"] = e.message
	res["file"] = e.file + ":" + strconv.Itoa(e.line)
	res["stack"] = e.stackTrace
	return res
}

func (e *logEntry) addSrcFileInfo() {
	var ok bool

	_, e.file, e.line, ok = runtime.Caller(2)
	if !ok {
		e.file = "???"
		e.line = 0
	}
}

func (e *logEntry) addStackTrace() { // this can be optimized in v2
	sbuff := make([]uintptr, maxStackLength)
	length := runtime.Callers(3, sbuff[:])
	stack := sbuff[:length]

	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		if strings.Contains(frame.File, "runtime/") {
			continue
		}

		e.stackTrace = append(e.stackTrace, fmt.Sprintf("%s:%d - %s", frame.File, frame.Line, frame.Function))

		if !more {
			break
		}
	}
}

func (e *logEntry) AddFields(f Fields) *logEntry {
	if e.fields == nil {
		e.fields = make(Fields)
	}

	for k, v := range f {
		e.fields[k] = v
	}

	return e
}

// NewError create and return *logEntry of error type.
// Arguments are handled in the manner of fmt.Print.
// Also adds file and line number.
func NewError(e interface{}) *logEntry {
	if v, ok := e.(*logEntry); ok {
		v.addSrcFileInfo()
		return v
	}

	var res *logEntry
	if v, ok := e.(error); ok {
		res = &logEntry{ltype: TError, message: v.Error(), err: v}
	} else if v, ok := e.(string); ok {
		res = &logEntry{ltype: TError, message: v}
	} else {
		res = &logEntry{ltype: TError, message: "Can't create new error!"}
	}

	res.addStackTrace()
	res.addSrcFileInfo()
	return res
}

// NewError create and return *logEntry of error type.
// Arguments are handled in the manner of fmt.Printf
// Also adds file and line number.
func NewErrorf(s string, v ...interface{}) (res *logEntry) {
	res = &logEntry{ltype: TError, message: fmt.Sprintf(s, v...)}
	res.addStackTrace()
	res.addSrcFileInfo()
	return
}

func NewWarning(s string) (res *logEntry) {
	res = &logEntry{ltype: TWarning, message: s}
	return
}

func NewWarningf(s string, v ...interface{}) (res *logEntry) {
	res = &logEntry{ltype: TWarning, message: fmt.Sprintf(s, v...)}
	return
}

func NewInfo(s string) (res *logEntry) {
	res = &logEntry{ltype: TInfo, message: s}
	return
}

func NewInfof(s string, v ...interface{}) (res *logEntry) {
	res = &logEntry{ltype: TInfo, message: fmt.Sprintf(s, v...)}
	return
}

func NewDebug(s string) (res *logEntry) {
	res = &logEntry{ltype: TDebug, message: s}
	res.addSrcFileInfo()
	return
}

func NewDebugf(s string, v ...interface{}) (res *logEntry) {
	res = &logEntry{ltype: TDebug, message: fmt.Sprintf(s, v...)}
	res.addSrcFileInfo()
	return
}
