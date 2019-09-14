package gctrl

import (
	"fmt"
	"reflect"
)

func errFunc(name string) func(interface{}) error {
	return func(err interface{}) error {
		return fmt.Errorf("gctrl.%s: %v", name, err)
	}
}

func parsingError(format string, args ...interface{}) (reflect.Value, []reflect.Value, error) {
	err := fmt.Errorf("parsing: "+format, args...)
	return reflect.ValueOf(nil), nil, err
}
