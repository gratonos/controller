package controller

import (
	"fmt"
	"io"
	"reflect"
)

func coutput(writer io.Writer, msg string, color colorID) {
	fmt.Fprintln(writer, colorize(msg, color))
}

func output(writer io.Writer, format string, args ...interface{}) {
	coutput(writer, fmt.Sprintf(format, args...), green)
}

func outputError(writer io.Writer, format string, args ...interface{}) {
	coutput(writer, fmt.Sprintf(format, args...), red)
}

func outputPrompt(writer io.Writer) {
	output(writer, "%s", `type '_list("")' to show registered function list`)
}

func outputCallResult(writer io.Writer, results []reflect.Value) {
	if len(results) == 0 {
		output(writer, "%s", "<void>")
		return
	}

	color := green
	if hasError(results) {
		color = red
	}

	for i, result := range results {
		typ := fmt.Sprint(result.Type())
		kind := fmt.Sprint(result.Kind())
		var output string
		if typ == kind {
			output = fmt.Sprintf("[%d] %v: %v", i, typ, result)
		} else {
			output = fmt.Sprintf("[%d] %v(%v): %v", i, typ, kind, result)
		}
		coutput(writer, output, color)
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
