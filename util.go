package controller

import (
	"errors"
	"fmt"
	"reflect"
)

func errFunc(name string) func(string) error {
	return func(err string) error {
		return fmt.Errorf("controller.%s: %s", name, err)
	}
}

func parseError(msg string) (fn reflect.Value, args []reflect.Value, err error) {
	return reflect.ValueOf(nil), nil, errors.New(msg)
}
