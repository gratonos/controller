package controller

import (
	"net"
	"os"
	"path/filepath"
)

const dirPerm = 0770

func (this *Controller) ServeUnix(path string) error {
	errFn := errFunc("ServeUnix")

	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return errFn(err)
	}
	if err := removeIfExists(path); err != nil {
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

func removeIfExists(path string) error {
	ok, err := fileExists(path)
	if err != nil {
		return err
	}
	if ok {
		return os.Remove(path)
	} else {
		return nil
	}
}

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}
