package main

import (
	"github.com/gratonos/controller"
)

func test() int {
	return 0
}
func main() {
	ctrl := controller.New()
	ctrl.MustRegister(test, "test", "test")
	ctrl.ServeUnix("/tmp/controller")
}
