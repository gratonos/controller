package controller

import (
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"
)

type command struct {
	name string
	desc string
}

func cprint(writer io.Writer, msg string, color colorID) {
	fmt.Fprintln(writer, colorize(msg, color))
}

func printMsg(writer io.Writer, format string, args ...interface{}) {
	cprint(writer, fmt.Sprintf(format, args...), green)
}

func printError(writer io.Writer, format string, args ...interface{}) {
	cprint(writer, fmt.Sprintf(format, args...), red)
}

func printPrompt(writer io.Writer) {
	printMsg(writer, "%s", "type '-help' to show builtin commands")
}

func printHelp(writer io.Writer) {
	commands := []command{
		{"-help", "show builtin commands"},
		{"-list [name]", "show registered functions (case insensitive, wildcard '*' supported)"},
		{"-prompt <on/off>", "set whether to print prompt while serving begins"},
	}

	tw := tabwriter.NewWriter(writer, 0, 0, 4, ' ', 0)
	for _, cmd := range commands {
		printMsg(tw, "%s\t%s", cmd.name, cmd.desc)
	}
	tw.Flush()
}

func printFuncList(writer io.Writer, metaList []*funcMeta) {
	tw := tabwriter.NewWriter(writer, 0, 0, 4, ' ', 0)
	for _, meta := range metaList {
		printMsg(tw, "%s\t%v\t// %s", meta.name, meta.fn.Type(), meta.desc)
	}
	tw.Flush()
}

func printCallResult(writer io.Writer, results []reflect.Value) {
	if len(results) == 0 {
		printMsg(writer, "%s", "<void>")
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
		cprint(writer, output, color)
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
