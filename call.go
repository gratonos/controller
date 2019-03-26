package controller

import (
	"fmt"
	"io"
	"reflect"
)

func (ctrl *Controller) handleFuncCall(wt io.Writer, cmd string) {
	fn, args, err := ctrl.parseExpr(cmd)
	if err != nil {
		printError(wt, err.Error())
		return
	}

	printResult(wt, fn.Call(args))
}

func printError(wt io.Writer, msg string) {
	fmt.Fprintln(wt, colorize(msg, red))
}

func printResult(wt io.Writer, results []reflect.Value) {
	color := green
	if hasError(results) {
		color = red
	}
	for i, result := range results {
		output := fmt.Sprintf("[%d] %v: %v", i, result.Type(), result)
		fmt.Fprintln(wt, colorize(output, color))
	}
	if len(results) == 0 {
		fmt.Fprintln(wt, colorize("<void>", green))
	}
}

func isError(value reflect.Value) bool {
	typ := value.Type()
	return typ.PkgPath() == "" && typ.Name() == "error"
}

func hasError(results []reflect.Value) bool {
	for _, result := range results {
		if isError(result) && !result.IsNil() {
			return true
		}
	}
	return false
}
