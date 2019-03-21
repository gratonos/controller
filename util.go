package controller

import (
	"fmt"
)

func genErrFunc(name string) func(string) error {
	return func(err string) error {
		return fmt.Errorf("controller.%s: %s", name, err)
	}
}
