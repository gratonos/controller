package controller

import (
	"errors"
	"fmt"
	"reflect"
)

func errFunc(name string) func(interface{}) error {
	return func(err interface{}) error {
		return fmt.Errorf("controller.%s: %v", name, err)
	}
}

func parseError(msg string) (fn reflect.Value, args []reflect.Value, err error) {
	return reflect.ValueOf(nil), nil, errors.New(msg)
}
