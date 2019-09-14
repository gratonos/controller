package main

import (
	"fmt"
	"reflect"

	"github.com/gratonos/gctrl"
)

func test() int {
	return 0
}

func main() {
	gctrl.MustRegister(test, "test", "test with unix domain socket")

	config := gctrl.ServeUnixConfig{
		OnConn: func(seq int64) {
			fmt.Printf("on connection(%d)\n", seq)
		},
		OnDisconn: func(seq int64, err error) {
			fmt.Printf("on disconnection(%d): err: %v\n", seq, err)
		},
		BeforeCall: func(connID int64, literal string) {
			fmt.Printf("call (conn: %d): %s\n", connID, literal)
		},
		AfterCall: func(connID int64, _ []reflect.Value, err error) {
			if err != nil {
				fmt.Printf("call (conn: %d) failed: %v\n", connID, err)
			} else {
				fmt.Printf("call (conn: %d) ok\n", connID)
			}
		},
	}

	panic("gctrl.ServeUnix: " + gctrl.ServeUnix("/tmp/gctrl", config).Error())
}
