package main

import (
	ctrl "github.com/gratonos/controller"
)

func test() int {
	return 0
}

func main() {
	ctrl.MustRegister(test, "test", "test with unix domain socket")
	ctrl.ServeUnix("/tmp/controller", ctrl.ServeUnixConfig{})
}
