package controller

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
)

type funcMeta struct {
	fn   reflect.Value
	in   []reflect.Type
	out  []reflect.Type
	name string
	desc string
}

type Controller struct {
	funcs  map[string]*funcMeta
	rwlock sync.RWMutex
}

const unixPerm = 0770

var defaultController = New()

func New() *Controller {
	return &Controller{
		funcs: make(map[string]*funcMeta),
	}
}

func (ctrl *Controller) Register(fn interface{}, name, desc string) error {
	if err := ctrl.register(fn, name, desc); err != nil {
		return errFunc("Register")(err)
	}
	return nil
}

func (ctrl *Controller) MustRegister(fn interface{}, name, desc string) {
	if err := ctrl.register(fn, name, desc); err != nil {
		panic(errFunc("MustRegister")(err))
	}
}

func (ctrl *Controller) Serve(rw io.ReadWriter) error {
	if err := ctrl.serve(rw); err != nil {
		return errFunc("Serve")(err)
	}
	return nil
}

func (ctrl *Controller) register(fn interface{}, name, desc string) error {
	if fn == nil {
		return errors.New("fn must not be nil")
	}
	if name == "" {
		return errors.New("name must not be empty")
	}
	meta, err := genFuncMeta(fn, name, desc)
	if err != nil {
		return err
	}

	ctrl.rwlock.Lock()
	defer ctrl.rwlock.Unlock()

	if _, ok := ctrl.funcs[name]; ok {
		return fmt.Errorf("name '%s' has been registered", name)
	}
	ctrl.funcs[name] = meta
	return nil
}

func (ctrl *Controller) serve(rw io.ReadWriter) error {
	prompt(rw)

	scanner := bufio.NewScanner(rw)
	for scanner.Scan() {
		cmd := strings.TrimSpace(scanner.Text())
		if cmd == "" {
			continue
		}

		ctrl.rwlock.RLock()

		if builtin(cmd) {
			ctrl.handleBuiltin(rw, cmd)
		} else {
			ctrl.handleFuncCall(rw, cmd)
		}

		ctrl.rwlock.RUnlock()
	}

	return scanner.Err()
}

func genFuncMeta(fn interface{}, name, desc string) (*funcMeta, error) {
	typ := reflect.TypeOf(fn)
	if typ.Kind() != reflect.Func {
		return nil, errors.New("fn is not a function")
	}

	var in []reflect.Type
	for i := 0; i < typ.NumIn(); i++ {
		inType := typ.In(i)
		if err := checkType(inType.Kind()); err != nil {
			return nil, err
		}
		in = append(in, inType)
	}

	var out []reflect.Type
	for i := 0; i < typ.NumOut(); i++ {
		outType := typ.Out(i)
		out = append(out, outType)
	}

	meta := &funcMeta{
		fn:   reflect.ValueOf(fn),
		in:   in,
		out:  out,
		name: name,
		desc: desc,
	}
	return meta, nil
}

func checkType(kind reflect.Kind) error {
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return nil
	default:
		return errors.New("supported parameter types are bool, integers, " +
			"floats, string and their compatible custom types")
	}
}
