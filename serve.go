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
		printPrompt(rw)
	}

	scanner := bufio.NewScanner(rw)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		this.rwlock.RLock()

		handleFuncCall(rw, input, this.funcMap)

		this.rwlock.RUnlock()
	}

	return scanner.Err()
}
