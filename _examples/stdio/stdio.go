package main

import (
	"errors"

	"github.com/gratonos/controller"
)

type coord struct {
	x, y int64
}

func sum(a, b int64) int64 {
	return a + b
}

func product(a, b int64) int64 {
	return a * b
}

func concat(s1, s2 string) string {
	return s1 + s2
}

func makeCoord(x, y int64) *coord {
	return &coord{x, y}
}

func makeError() error {
	return errors.New("error: test")
}

func main() {
	ctrl := controller.New()
	ctrl.MustRegister(sum, "sum", "sum of two int64")
	ctrl.MustRegister(product, "product", "product of two int64")
	ctrl.MustRegister(concat, "concat", "concatenate two string")
	ctrl.MustRegister(makeCoord, "makeCoord", "make a coordinate")
	ctrl.MustRegister(makeError, "makeError", "make an error")
	ctrl.Serve(controller.Stdio())
}
