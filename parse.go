package controller

import (
	"errors"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

type numberType int

const (
	typeInt numberType = iota
	typeUint
	typeFloat
)

var parserMap = map[reflect.Kind]func(ast.Expr) (reflect.Value, error){
	reflect.Bool:    parseBool,
	reflect.Int64:   parseInt64,
	reflect.Uint64:  parseUint64,
	reflect.Float64: parseFloat64,
	reflect.String:  parseString,
}

func parseError(err error) (reflect.Value, error) {
	return reflect.ValueOf(nil), err
}

func parseBool(expr ast.Expr) (reflect.Value, error) {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return parseError(errors.New("bool !ok"))
	}
	if ident.Name == "true" {
		return reflect.ValueOf(true), nil
	} else if ident.Name == "false" {
		return reflect.ValueOf(false), nil
	} else {
		return parseError(errors.New("bool invalid"))
	}
}

func parseInt64(expr ast.Expr) (reflect.Value, error) {
	switch expr := expr.(type) {
	case *ast.BasicLit:
		return parseInt(expr, true)
	case *ast.UnaryExpr:
		return parseUnary(expr, typeInt)
	default:
		return parseError(errors.New("int64 !ok"))
	}
}

func parseUint64(expr ast.Expr) (reflect.Value, error) {
	switch expr := expr.(type) {
	case *ast.BasicLit:
		return parseUint(expr)
	case *ast.UnaryExpr:
		return parseUnary(expr, typeUint)
	default:
		return parseError(errors.New("uint64 !ok"))
	}
}

func parseFloat64(expr ast.Expr) (reflect.Value, error) {
	switch expr := expr.(type) {
	case *ast.BasicLit:
		return parseFloat(expr, true)
	case *ast.UnaryExpr:
		return parseUnary(expr, typeFloat)
	default:
		return parseError(errors.New("float64 !ok"))
	}
}

func parseString(expr ast.Expr) (reflect.Value, error) {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return parseError(errors.New("string !ok"))
	}
	if lit.Kind != token.STRING {
		return parseError(errors.New("string invalid"))
	}
	str, err := strconv.Unquote(lit.Value)
	if err != nil {
		return parseError(errors.New("Unquote error"))
	}
	return reflect.ValueOf(str), nil
}

func parseUnary(expr *ast.UnaryExpr, numType numberType) (reflect.Value, error) {
	x, ok := expr.X.(*ast.BasicLit)
	if !ok {
		return parseError(errors.New("unary invalid"))
	}
	switch expr.Op {
	case token.ADD:
		switch numType {
		case typeInt:
			return parseInt(x, true)
		case typeUint:
			return parseUint(x)
		case typeFloat:
			return parseFloat(x, true)
		default:
			panic("impossible")
		}
	case token.SUB:
		switch numType {
		case typeInt:
			return parseInt(x, false)
		case typeUint:
			return parseError(errors.New("neg uint"))
		case typeFloat:
			return parseFloat(x, false)
		default:
			panic("impossible")
		}
	default:
		return parseError(errors.New("unary !ok"))
	}
}

func parseInt(expr *ast.BasicLit, positive bool) (reflect.Value, error) {
	if expr.Kind != token.INT {
		return parseError(errors.New("int !ok"))
	}
	value, err := strconv.ParseInt(expr.Value, 0, 64)
	if err != nil {
		return parseError(errors.New("int invalid"))
	}
	if !positive {
		value = -value
	}
	return reflect.ValueOf(value), nil
}

func parseUint(expr *ast.BasicLit) (reflect.Value, error) {
	if expr.Kind != token.INT {
		return parseError(errors.New("uint !ok"))
	}
	value, err := strconv.ParseUint(expr.Value, 0, 64)
	if err != nil {
		return parseError(errors.New("uint invalid"))
	}
	return reflect.ValueOf(value), nil
}

func parseFloat(expr *ast.BasicLit, positive bool) (reflect.Value, error) {
	if expr.Kind != token.FLOAT && expr.Kind != token.INT {
		return parseError(errors.New("float !ok"))
	}
	value, err := strconv.ParseFloat(expr.Value, 64)
	if err != nil {
		return parseError(errors.New("float invalid"))
	}
	if !positive {
		value = -value
	}
	return reflect.ValueOf(value), nil
}
