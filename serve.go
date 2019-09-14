package gctrl

import (
	"bufio"
	"io"
	"strings"
)

type ServeConfig struct {
	NoPrompt   bool
	NoColoring bool
}

func (this *Controller) Serve(rw io.ReadWriter, config ServeConfig) error {
	if err := this.serve(rw, config); err != nil {
		return errFunc("Serve")(err)
	}
	return nil
}

func (this *Controller) serve(rw io.ReadWriter, config ServeConfig) error {
	output := &output{
		Writer:   rw,
		Coloring: !config.NoColoring,
	}

	if !config.NoPrompt {
		output.Prompt()
	}

	scanner := bufio.NewScanner(rw)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input != "" {
			this.doCall(input, output)
		}
	}

	return scanner.Err()
}

func (this *Controller) doCall(literal string, output *output) {
	this.rwlock.RLock()
	defer this.rwlock.RUnlock()

	results, err := call(literal, this.funcMap)
	if err != nil {
		output.Error("%v", err)
		return
	}

	output.CallResult(results)
}
