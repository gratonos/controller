package controller

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
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

func New() *Controller {
	return &Controller{
		funcs: make(map[string]*funcMeta),
	}
}

func (this *Controller) Register(fn interface{}, name, desc string) error {
	if err := this.register(fn, name, desc); err != nil {
		return errFunc("Register")(err)
	}
	return nil
}

func (this *Controller) MustRegister(fn interface{}, name, desc string) {
	if err := this.register(fn, name, desc); err != nil {
		panic(errFunc("MustRegister")(err))
	}
}

func (this *Controller) Serve(rw io.ReadWriter) error {
	if err := this.serve(rw); err != nil {
		return errFunc("Serve")(err)
	}
	return nil
}

func (this *Controller) register(fn interface{}, name, desc string) error {
	if fn == nil {
		return errors.New("fn must not be nil")
	}
	if name == "" {
		return errors.New("name must not be empty")
	}
	expr, err := parser.ParseExpr(name)
	if err != nil {
		return errors.New("name is invalid")
	}
	_, ok := expr.(*ast.Ident)
	if !ok {
		return errors.New("name is not a valid identity")
	}
	meta, err := makeFuncMeta(fn, name, desc)
	if err != nil {
		return err
	}

	this.rwlock.Lock()
	defer this.rwlock.Unlock()

	if _, ok := this.funcs[name]; ok {
		return fmt.Errorf("name '%s' has been registered", name)
	}
	this.funcs[name] = meta

	return nil
}

func (this *Controller) serve(rw io.ReadWriter) error {
	prompt(rw)

	scanner := bufio.NewScanner(rw)
	for scanner.Scan() {
		cmd := strings.TrimSpace(scanner.Text())
		if cmd == "" {
			continue
		}

		this.rwlock.RLock()

		if builtin(cmd) {
			this.handleBuiltin(rw, cmd)
		} else {
			this.handleFuncCall(rw, cmd)
		}

		this.rwlock.RUnlock()
	}

	return scanner.Err()
}

func makeFuncMeta(fn interface{}, name, desc string) (*funcMeta, error) {
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
			"floats, string and their type aliases")
	}
}
