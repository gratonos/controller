package controller

import (
	"fmt"
	"go/ast"
	"go/parser"
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

func (ctrl *Controller) parseExpr(cmd string) (fn reflect.Value, args []reflect.Value, err error) {
	expr, err := parser.ParseExpr(cmd)
	if err != nil {
		return parseError(fmt.Sprintf("parse: '%s' is not a valid expression", cmd))
	}
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return parseError(fmt.Sprintf("parse: '%s' is not a function call expression", cmd))
	}
	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return parseError(fmt.Sprintf("parse: '%s' is not a valid function call expression", cmd))
	}
	meta, ok := ctrl.funcs[ident.Name]
	if !ok {
		return parseError(fmt.Sprintf("parse: function '%s' is not registered", ident.Name))
	}
	if len(meta.in) != len(call.Args) {
		return parseError(fmt.Sprintf("parse: unmatched argument count, want %d, have %d",
			len(meta.in), len(call.Args)))
	}

	for i := 0; i < len(meta.in); i++ {
		text := cmd[call.Args[i].Pos()-1 : call.Args[i].End()-1]
		arg, err := parseArg(text, meta.in[i])
		if err != nil {
			return parseError(fmt.Sprintf("parse: argument[%d] '%s' is not a valid "+
				"literal of type %v", i, text, meta.in[i]))
		}
		args = append(args, arg)
	}

	return meta.fn, args, nil
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
