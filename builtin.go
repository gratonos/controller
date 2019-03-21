package controller

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

func prompt(wt io.Writer) {
	fmt.Fprintln(wt, "conn")
}

func builtin(cmd string) bool {
	return strings.HasPrefix(cmd, "-")
}

func (ctrl *Controller) handleBuiltin(cmd string, wt io.Writer) {
	switch cmd {
	case "-list":
		ctrl.handleList(wt)
	default:
		// todo
	}
}

func (ctrl *Controller) handleList(wt io.Writer) {
	var funcList []string
	for name := range ctrl.funcs {
		funcList = append(funcList, name)
	}
	sort.Strings(funcList)

	for _, name := range funcList {
		meta := ctrl.funcs[name]
		fmt.Fprintf(wt, "%s: %s: %s\n", meta.name, meta.fn.Type(), meta.desc)
	}
}
