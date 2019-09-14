package gctrl

import (
	"io"
)

var defaultController *Controller

func init() {
	defaultController = New()
}

func Register(fn interface{}, name, desc string) error {
	return defaultController.Register(fn, name, desc)
}

func MustRegister(fn interface{}, name, desc string) {
	defaultController.MustRegister(fn, name, desc)
}

func Call(text string) ([]interface{}, error) {
	return defaultController.Call(text)
}

func Serve(rw io.ReadWriter, config ServeConfig) error {
	return defaultController.Serve(rw, config)
}

func ServeUnix(path string, config ServeUnixConfig) error {
	return defaultController.ServeUnix(path, config)
}
