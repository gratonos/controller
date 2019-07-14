package controller

import (
	"fmt"
	"go/ast"
	"go/parser"
	"io"
	"reflect"
	"strconv"
)

func (this *Controller) handleFuncCall(writer io.Writer, cmd string) {
	fn, args, err := this.parseExpr(cmd)
	if err != nil {
		printError(writer, "%v", err)
		return
	}

	printCallResult(writer, fn.Call(args))
}

func (this *Controller) parseExpr(cmd string) (fn reflect.Value, args []reflect.Value, err error) {
	expr, err := parser.ParseExpr(cmd)
	if err != nil {
		return parseError("'%s' is not a valid expression", cmd)
	}
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return parseError("'%s' is not a function call expression", cmd)
	}
	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return parseError("'%s' is not a valid function call expression", cmd)
	}
	meta, ok := this.funcs[ident.Name]
	if !ok {
		return parseError("function '%s' is not registered", ident.Name)
	}
	if len(meta.in) != len(call.Args) {
		return parseError("unmatched argument count, want %d, have %d",
			len(meta.in), len(call.Args))
	}

	for i := 0; i < len(meta.in); i++ {
		text := cmd[call.Args[i].Pos()-1 : call.Args[i].End()-1]
		arg, err := parseArg(text, meta.in[i])
		if err != nil {
			return parseError("argument[%d] '%s' is not a valid literal of type %v",
				i, text, meta.in[i])
		}
		args = append(args, arg)
	}

	return meta.fn, args, nil
}

func parseArg(text string, expectedType reflect.Type) (reflect.Value, error) {
	var value interface{}
	var err error
	switch expectedType.Kind() {
	case reflect.Bool:
		value, err = parseBool(text)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err = parseInt(text, int(expectedType.Size()*8))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err = parseUint(text, int(expectedType.Size()*8))
	case reflect.Float32, reflect.Float64:
		value, err = parseFloat(text, int(expectedType.Size()*8))
	case reflect.String:
		value, err = parseString(text)
	default:
		panic("impossible")
	}

	if err != nil {
		return reflect.ValueOf(nil), err
	}
	return reflect.ValueOf(value).Convert(expectedType), nil
}

func parseBool(text string) (interface{}, error) {
	switch text {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return nil, fmt.Errorf("'%s' is not a bool literal")
	}
}

func parseInt(text string, bitSize int) (interface{}, error) {
	return strconv.ParseInt(text, 0, bitSize)
}

func parseUint(text string, bitSize int) (interface{}, error) {
	return strconv.ParseUint(text, 0, bitSize)
}

func parseFloat(text string, bitSize int) (interface{}, error) {
	return strconv.ParseFloat(text, bitSize)
}

func parseString(text string) (interface{}, error) {
	return strconv.Unquote(text)
}
