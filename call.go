package controller

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"io"
	"reflect"
)

func (ctrl *Controller) handleFuncCall(cmd string, wt io.Writer) {
	expr, err := parser.ParseExpr(cmd)
	if err != nil {
		printError(wt, err)
		return
	}

	call, ok := expr.(*ast.CallExpr)
	if !ok {
		printError(wt, errors.New("call !ok"))
		return
	}

	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		printError(wt, errors.New("ident !ok"))
		return
	}

	meta, ok := ctrl.funcs[ident.Name]
	if !ok {
		printError(wt, errors.New("meta !ok"))
		return
	}

	if len(meta.in) != len(call.Args) {
		printError(wt, errors.New("in !ok"))
		return
	}

	var args []reflect.Value
	for i := 0; i < len(meta.in); i++ {
		arg, err := parserMap[meta.in[i].Kind()](call.Args[i])
		if err != nil {
			printError(wt, err)
			return
		}
		args = append(args, arg)
	}

	printResult(wt, meta.fn.Call(args))
}

func printError(wt io.Writer, err error) {
	fmt.Fprintln(wt, err)
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
