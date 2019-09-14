package gctrl

import (
	"net"
	"os"
	"path/filepath"
	"reflect"

	gos "github.com/gratonos/goutil/os"
)

const dirPerm = 0770

type OnConn func(connID int64)
type OnDisconn func(connID int64, err error)
type UnixBeforeCall func(connID int64, literal string)
type UnixAfterCall func(connID int64, results []reflect.Value, err error)

type ServeUnixConfig struct {
	NoPrompt   bool
	NoColoring bool
	OnConn     OnConn
	OnDisconn  OnDisconn
	BeforeCall UnixBeforeCall
	AfterCall  UnixAfterCall
}

func (this *Controller) ServeUnix(path string, config ServeUnixConfig) error {
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

	var serial int64 = 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			return errFn(err)
		}

		go func(connID int64) {
			if config.OnConn != nil {
				config.OnConn(connID)
			}

			err := this.serve(conn, makeServeConfig(connID, &config))
			conn.Close()

			if config.OnDisconn != nil {
				config.OnDisconn(connID, err)
			}
		}(serial)

		serial++
	}
}

func makeServeConfig(connID int64, config *ServeUnixConfig) *ServeConfig {
	var before BeforeCall
	if config.BeforeCall != nil {
		before = func(literal string) {
			config.BeforeCall(connID, literal)
		}
	}

	var after AfterCall
	if config.AfterCall != nil {
		after = func(results []reflect.Value, err error) {
			config.AfterCall(connID, results, err)
		}
	}

	return &ServeConfig{
		NoPrompt:   config.NoPrompt,
		NoColoring: config.NoColoring,
		BeforeCall: before,
		AfterCall:  after,
	}
}
