package main

import (
	"math/rand"
	"time"

	ctrl "github.com/gratonos/controller"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ID int

type Coord struct {
	x, y int
}

type Error struct {
	msg string
}

func (err *Error) Error() string {
	return err.msg
}

func void() {}

func and(a, b bool) bool {
	return a && b
}

func addi(a, b int) int {
	return a + b
}

func addi8(a, b int8) int8 {
	return a + b
}

func addi16(a, b int16) int16 {
	return a + b
}

func addi32(a, b int32) int32 {
	return a + b
}

func addi64(a, b int64) int64 {
	return a + b
}

func addu(a, b uint) uint {
	return a + b
}

func addu8(a, b uint8) uint8 {
	return a + b
}

func addu16(a, b uint16) uint16 {
	return a + b
}

func addu32(a, b uint32) uint32 {
	return a + b
}

func addu64(a, b uint64) uint64 {
	return a + b
}

func addf32(a, b float32) float32 {
	return a + b
}

func addf64(a, b float64) float64 {
	return a + b
}

func concat(s1, s2 string) string {
	return s1 + s2
}

func makeCoord(x, y int) *Coord {
	return &Coord{x, y}
}

func alias(id ID) ID {
	return id
}

func randError() (int, error) {
	randn := rand.Int()
	if randn%2 == 0 {
		return randn, nil
	} else {
		return randn, &Error{"odd number"}
	}
}

func main() {
	ctrl.MustRegister(void, "void", "do nothing")

	ctrl.MustRegister(and, "and", "logic AND")

	ctrl.MustRegister(addi, "addi", "sum two ints")
	ctrl.MustRegister(addi8, "addi8", "sum two int8s")
	ctrl.MustRegister(addi16, "addi16", "sum two int16s")
	ctrl.MustRegister(addi32, "addi32", "sum two int32s")
	ctrl.MustRegister(addi64, "addi64", "sum two int64s")

	ctrl.MustRegister(addu, "addu", "sum two uints")
	ctrl.MustRegister(addu8, "addu8", "sum two uint8s")
	ctrl.MustRegister(addu16, "addu16", "sum two uint16s")
	ctrl.MustRegister(addu32, "addu32", "sum two uint32s")
	ctrl.MustRegister(addu64, "addu64", "sum two uint64s")

	ctrl.MustRegister(addf32, "addf32", "sum two float32s")
	ctrl.MustRegister(addf64, "addf64", "sum two float64s")

	ctrl.MustRegister(concat, "concat", "concatenate two strings")

	ctrl.MustRegister(makeCoord, "makeCoord", "make a coordinate")
	ctrl.MustRegister(alias, "alias", "alias of ID")

	ctrl.MustRegister(randError, "randError", "rand of error")

	ctrl.Serve(ctrl.Stdio())
}
