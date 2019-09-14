package gctrl

import (
	"bufio"
	"io"
	"reflect"
	"strings"
)

type BeforeCall func(literal string)
type AfterCall func(results []reflect.Value, err error)

type ServeConfig struct {
	NoPrompt   bool
	NoColoring bool
	BeforeCall BeforeCall
	AfterCall  AfterCall
}

func (this *Controller) Serve(rw io.ReadWriter, config ServeConfig) error {
	if err := this.serve(rw, &config); err != nil {
		return errFunc("Serve")(err)
	}
	return nil
}

func (this *Controller) serve(rw io.ReadWriter, config *ServeConfig) error {
	printer := &printer{
		Writer:   rw,
		Coloring: !config.NoColoring,
	}

	if !config.NoPrompt {
		printer.Prompt()
	}

	scanner := bufio.NewScanner(rw)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input != "" {
			if config.BeforeCall != nil {
				config.BeforeCall(input)
			}

			results, err := this.doCall(input)
			if err != nil {
				printer.Error("%v", err)
			} else {
				printer.Result(results)
			}

			if config.AfterCall != nil {
				config.AfterCall(results, err)
			}
		}
	}

	return scanner.Err()
}

func (this *Controller) doCall(literal string) ([]reflect.Value, error) {
	this.rwlock.RLock()
	defer this.rwlock.RUnlock()

	return call(literal, this.funcMap)
}
