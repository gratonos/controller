package controller

import (
	"net"
	"os"
	"path/filepath"
)

const dirPerm = 0770

func (ctrl *Controller) ServeUnix(path string) error {
	errfn := errFunc("ServeUnix")

	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return errfn(err)
	}
	if err := checkAndRemove(path); err != nil {
		return errfn(err)
	}

	listener, err := net.Listen("unix", path)
	if err != nil {
		return errfn(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return errfn(err)
		}
		go func() {
			ctrl.serve(conn)
			conn.Close()
		}()
	}
}

func checkAndRemove(path string) error {
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	return os.Remove(path)
}
