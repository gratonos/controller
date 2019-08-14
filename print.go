package controller

import (
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"
)

func prompt(writer io.Writer) {
	msg := "type '-list [name]' to show registered functions " +
		"(case insensitive, wildcard '*' supported)"
	fmt.Fprintln(writer, colorize(msg, green))
}

func printError(writer io.Writer, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(writer, colorize(msg, red))
}

func printFuncList(writer io.Writer, metaList []*funcMeta) {
	tw := tabwriter.NewWriter(writer, 0, 0, 4, ' ', 0)

	for _, meta := range metaList {
		output := fmt.Sprintf("%s\t%v\t// %s", meta.name, meta.fn.Type(), meta.desc)
		fmt.Fprintln(tw, colorize(output, green))
	}

	tw.Flush()
}

func printCallResult(writer io.Writer, results []reflect.Value) {
	if len(results) == 0 {
		fmt.Fprintln(writer, colorize("<void>", green))
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
		fmt.Fprintln(writer, colorize(output, color))
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
