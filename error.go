package controller

import (
	"fmt"
	"reflect"
)

func errFunc(name string) func(interface{}) error {
	return func(err interface{}) error {
		return fmt.Errorf("controller.%s: %v", name, err)
	}
}

func parseError(format string, args ...interface{}) (reflect.Value, []reflect.Value, error) {
	err := fmt.Errorf("parse: "+format, args...)
	return reflect.ValueOf(nil), nil, err
}
