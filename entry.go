package log

import (
	"fmt"
	"runtime"
	"strconv"
)

const (
	TError   = "error"
	TWarning = "warning"
	TInfo    = "info"
	TDebug   = "debug"
)

type Entrier interface {
	ToFields() Fields
	Parent() error
}

type Fields map[string]interface{}

type logEntry struct {
	ltype   string
	message string
	line    int
	file    string
	fields  Fields

	parent error
}

func (e *logEntry) Error() string {
	if e.ltype == "error" {
		return e.message
	}

	return ""
}

func (e *logEntry) Parent() error {
	return e.parent
}

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
		res = &logEntry{ltype: TError, message: v.Error(), parent: v}
	} else if v, ok := e.(string); ok {
		res = &logEntry{ltype: TError, message: v}
	} else {
		res = &logEntry{ltype: TError, message: "Can't create new error!"}
	}

	res.addSrcFileInfo()
	return res
}

// NewError create and return *logEntry of error type.
// Arguments are handled in the manner of fmt.Printf
// Also adds file and line number.
func NewErrorf(s string, v ...interface{}) (res *logEntry) {
	res = &logEntry{ltype: TError, message: fmt.Sprintf(s, v...)}
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
