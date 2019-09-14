package main

import (
	"github.com/gratonos/gctrl"
)

func test() int {
	return 0
}

func main() {
	gctrl.MustRegister(test, "test", "test with unix domain socket")
	gctrl.ServeUnix("/tmp/gctrl", gctrl.ServeUnixConfig{})
}
