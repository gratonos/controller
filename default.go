package controller

import (
	"io"
)

var defaultController *Controller

func init() {
	defaultController = New(true)
}

func Prompt() bool {
	return defaultController.Prompt()
}

func SetPrompt(prompt bool) {
	defaultController.SetPrompt(prompt)
}

func Register(fn interface{}, name, desc string) error {
	return defaultController.Register(fn, name, desc)
}

func MustRegister(fn interface{}, name, desc string) {
	defaultController.MustRegister(fn, name, desc)
}

func Serve(rw io.ReadWriter) error {
	return defaultController.Serve(rw)
}

func ServeUnix(path string) error {
	return defaultController.ServeUnix(path)
}
