package controller

import (
	"fmt"
	"go/ast"
	"go/parser"
	"io"
	"reflect"
)

func (ctrl *Controller) handleFuncCall(cmd string, wt io.Writer) {
	expr, err := parser.ParseExpr(cmd)
	if err != nil {
		printError(wt, err.Error())
		return
	}

	call, ok := expr.(*ast.CallExpr)
	if !ok {
		printError(wt, "call !ok")
		return
	}

	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		printError(wt, "ident !ok")
		return
	}

	meta, ok := ctrl.funcs[ident.Name]
	if !ok {
		printError(wt, "meta !ok")
		return
	}

	if len(meta.in) != len(call.Args) {
		printError(wt, "in !ok")
		return
	}

	var args []reflect.Value
	for i := 0; i < len(meta.in); i++ {
		arg, err := parserMap[meta.in[i].Kind()](call.Args[i])
		if err != nil {
			printError(wt, err.Error())
			return
		}
		args = append(args, arg)
	}

	printResult(wt, meta.fn.Call(args))
}

func printError(wt io.Writer, msg string) {
	fmt.Fprintln(wt, msg)
}

func printResult(wt io.Writer, results []reflect.Value) {
	for _, result := range results {
		fmt.Fprintln(wt, result)
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
