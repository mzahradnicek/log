package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Logger struct {
	mu  sync.Mutex
	out io.Writer
	buf []byte
}

func (l *Logger) Save(e interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var f Fields
	if v, ok := e.(Entrier); ok {
		f = v.ToFields()
	} else if v, ok := e.(error); ok {
		f = Fields{"type": TError, "msg": v.Error()}
	} else {
		fmt.Println("notype")
		return
	}

	f["time"] = time.Now()

	var err error
	l.buf, err = json.Marshal(f)
	if err != nil {
		l.buf = []byte("{ type=\"logerror\" msg=\"" + err.Error() + "\"}")
		return
	}

	_, err = l.out.Write(l.buf)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}

func (l *Logger) Error(s string) {
	e := &logEntry{ltype: TError, message: s}
	e.addSrcFileInfo()
	l.Save(e)
}

func (l *Logger) Errorf(s string, v ...interface{}) {
	e := &logEntry{ltype: TError, message: fmt.Sprintf(s, v...)}
	e.addSrcFileInfo()
	l.Save(e)
}

func (l *Logger) Debug(s string) {
	e := &logEntry{ltype: TDebug, message: s}
	e.addSrcFileInfo()
	l.Save(e)
}

func (l *Logger) Debugf(s string, v ...interface{}) {
	e := &logEntry{ltype: TDebug, message: fmt.Sprintf(s, v...)}
	e.addSrcFileInfo()
	l.Save(e)
}

func (l *Logger) Info(s string) {
	e := &logEntry{ltype: TInfo, message: s}
	l.Save(e)
}

func (l *Logger) Infof(s string, v ...interface{}) {
	e := &logEntry{ltype: TInfo, message: fmt.Sprintf(s, v...)}
	l.Save(e)
}

func (l *Logger) Warning(s string) {
	e := &logEntry{ltype: TWarning, message: s}
	l.Save(e)
}

func (l *Logger) Warningf(s string, v ...interface{}) {
	e := &logEntry{ltype: TWarning, message: fmt.Sprintf(s, v...)}
	l.Save(e)
}

func (l *Logger) SetOutput(o io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = o
}

func New(out io.Writer) *Logger { // writer
	return &Logger{out: out}
}

var std = New(os.Stderr)

func SetOutput(o io.Writer) {
	std.SetOutput(o)
}

func Save(e interface{}) {
	std.Save(e)
}

func Error(s string) {
	e := &logEntry{ltype: TError, message: s}
	e.addSrcFileInfo()
	std.Save(e)
}

func Errorf(s string, v ...interface{}) {
	e := &logEntry{ltype: TError, message: fmt.Sprintf(s, v...)}
	e.addSrcFileInfo()
	std.Save(e)
}

func Debug(s string) {
	e := &logEntry{ltype: TDebug, message: s}
	e.addSrcFileInfo()
	std.Save(e)
}

func Debugf(s string, v ...interface{}) {
	e := &logEntry{ltype: TDebug, message: fmt.Sprintf(s, v...)}
	e.addSrcFileInfo()
	std.Save(e)
}

func Info(s string) {
	e := &logEntry{ltype: TInfo, message: s}
	std.Save(e)
}

func Infof(s string, v ...interface{}) {
	e := &logEntry{ltype: TInfo, message: fmt.Sprintf(s, v...)}
	std.Save(e)
}

func Warning(s string) {
	e := &logEntry{ltype: TWarning, message: s}
	std.Save(e)
}

func Warningf(s string, v ...interface{}) {
	e := &logEntry{ltype: TWarning, message: fmt.Sprintf(s, v...)}
	std.Save(e)
}
