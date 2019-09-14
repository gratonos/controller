package gctrl

import (
	"io"
	"os"
)

type stdio struct{}

func (stdio) Read(bs []byte) (n int, err error) {
	return os.Stdin.Read(bs)
}

func (stdio) Write(bs []byte) (n int, err error) {
	return os.Stdout.Write(bs)
}

func Stdio() io.ReadWriter {
	return stdio{}
}
