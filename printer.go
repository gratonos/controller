package gctrl

import (
	"fmt"
	"io"
	"reflect"
)

type colorID int

const (
	black colorID = iota + 30
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

type printer struct {
	Writer   io.Writer
	Coloring bool
}

func (this *printer) Normal(format string, args ...interface{}) {
	this.write(green, format, args...)
}

func (this *printer) Error(format string, args ...interface{}) {
	this.write(red, format, args...)
}

func (this *printer) Prompt() {
	this.Normal("%s", `type '_list("")' to show registered function list`)
}

func (this *printer) Result(results []reflect.Value) {
	if len(results) == 0 {
		this.Normal("%s", "<void>")
		return
	}

	printer := this.Normal
	if hasError(results) {
		printer = this.Error
	}

	for i, result := range results {
		typ := fmt.Sprint(result.Type())
		kind := fmt.Sprint(result.Kind())
		if typ == kind {
			printer("[%d] %v: %v", i, typ, result)
		} else {
			printer("[%d] %v(%v): %v", i, typ, kind, result)
		}
	}
}

func (this *printer) write(color colorID, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if this.Coloring {
		fmt.Fprintln(this.Writer, colorize(msg, color))
	} else {
		fmt.Fprintln(this.Writer, msg)
	}
}

func hasError(results []reflect.Value) bool {
	for _, result := range results {
		if isError(result) && !result.IsNil() {
			return true
		}
	}
	return false
}

func isError(value reflect.Value) bool {
	typ := value.Type()
	return typ.PkgPath() == "" && typ.Name() == "error"
}

func colorize(text string, color colorID) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", color, text)
}
