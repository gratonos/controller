package controller

import (
	"bufio"
	"io"
	"strings"
)

type ServeConfig struct {
	noPrompt bool
}

func (this *Controller) Serve(rw io.ReadWriter, config ServeConfig) error {
	if err := this.serve(rw, config); err != nil {
		return errFunc("Serve")(err)
	}
	return nil
}

func (this *Controller) serve(rw io.ReadWriter, config ServeConfig) error {
	if !config.noPrompt {
		outputPrompt(rw)
	}

	scanner := bufio.NewScanner(rw)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input != "" {
			this.doCall(rw, input)
		}
	}

	return scanner.Err()
}

func (this *Controller) doCall(writer io.Writer, literal string) {
	this.rwlock.RLock()
	defer this.rwlock.RUnlock()

	results, err := call(literal, this.funcMap)
	if err != nil {
		outputError(writer, "%v", err)
		return
	}

	outputCallResult(writer, results)
}
