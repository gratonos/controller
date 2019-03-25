package controller

import (
	"fmt"
	"reflect"
	"strconv"
)

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
