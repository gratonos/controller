package controller

import (
	"fmt"
	"go/ast"
	"go/parser"
	"reflect"
	"strconv"
)

func call(literal string, funcMap map[string]*funcMeta) ([]reflect.Value, error) {
	fn, args, err := parse(literal, funcMap)
	if err != nil {
		return nil, err
	}
	return fn.Call(args), nil
}

func parse(literal string, funcMap map[string]*funcMeta) (
	fn reflect.Value, args []reflect.Value, err error) {

	call, err := parseExpr(literal)
	if err != nil {
		return parsingError("%v", err)
	}

	name, err := parseName(call, literal)
	if err != nil {
		return parsingError("%v", err)
	}

	meta, ok := funcMap[name]
	if !ok {
		return parsingError("function '%s' is not registered", name)
	}

	args, err = parseArgs(call, literal, meta)
	if err != nil {
		return parsingError("%v", err)
	}

	return meta.fn, args, nil
}

func parseExpr(literal string) (*ast.CallExpr, error) {
	expr, err := parser.ParseExpr(literal)
	if err != nil {
		return nil, fmt.Errorf("'%s' is not a valid expression", literal)
	}
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return nil, fmt.Errorf("'%s' is not a call expression", literal)
	}
	return call, nil
}

func parseName(call *ast.CallExpr, literal string) (string, error) {
	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		text := exprText(call.Fun, literal)
		return "", fmt.Errorf("'%s' is not a valid function name", text)
	}
	return ident.Name, nil
}

func parseArgs(call *ast.CallExpr, literal string, meta *funcMeta) ([]reflect.Value, error) {
	if len(meta.in) != len(call.Args) {
		return nil, fmt.Errorf("unmatched argument count, want %d, have %d",
			len(meta.in), len(call.Args))
	}

	var args []reflect.Value
	for i := 0; i < len(meta.in); i++ {
		text := exprText(call.Args[i], literal)
		arg, err := parseArg(text, meta.in[i])
		if err != nil {
			return nil, fmt.Errorf("argument[%d] '%s' is not a valid literal of type %v",
				i, text, meta.in[i])
		}
		args = append(args, arg)
	}

	return args, nil
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
		panic("checking argument failure")
	}

	if err != nil {
		return reflect.ValueOf(nil), err
	}
	return reflect.ValueOf(value).Convert(expectedType), nil
}

func parseBool(text string) (bool, error) {
	switch text {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("'%s' is not a bool literal")
	}
}

func parseInt(text string, bitSize int) (int64, error) {
	return strconv.ParseInt(text, 0, bitSize)
}

func parseUint(text string, bitSize int) (uint64, error) {
	return strconv.ParseUint(text, 0, bitSize)
}

func parseFloat(text string, bitSize int) (float64, error) {
	return strconv.ParseFloat(text, bitSize)
}

func parseString(text string) (string, error) {
	return strconv.Unquote(text)
}

func exprText(expr ast.Expr, text string) string {
	return text[expr.Pos()-1 : expr.End()-1]
}
