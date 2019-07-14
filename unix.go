package controller

import (
	"net"
	"os"
	"path/filepath"

	gos "github.com/gratonos/goutil/os"
)

const dirPerm = 0770

func (this *Controller) ServeUnix(path string) error {
	errFn := errFunc("ServeUnix")

	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return errFn(err)
	}
	if err := gos.RemoveIfExists(path); err != nil {
		return errFn(err)
	}

	listener, err := net.Listen("unix", path)
	if err != nil {
		return errFn(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return errFn(err)
		}
		go func() {
			_ = this.serve(conn)
			conn.Close()
		}()
	}
}
