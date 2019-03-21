package controller

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
)

func prompt(wt io.Writer) {
	fmt.Fprintln(wt, "type '-list [name]' to show registered functions "+
		"(case insensitive, wildcard '*' supported)")
}

func builtin(cmd string) bool {
	return strings.HasPrefix(cmd, "-")
}

func (ctrl *Controller) handleBuiltin(cmd string, wt io.Writer) {
	switch {
	case strings.HasPrefix(cmd, "-list"):
		ctrl.handleList(wt, strings.Fields(cmd)[1:])
	default:
		printError(wt, "unsupported command")
	}
}

func (ctrl *Controller) handleList(wt io.Writer, args []string) {
	argc := len(args)
	if argc > 1 {
		printError(wt, "-list: too many arguments")
		return
	}

	exp := ""
	if argc == 1 {
		exp = "(?i)" + strings.Replace(args[0], "*", ".*", -1)
	}
	reg, err := regexp.Compile(exp)
	if err != nil {
		printError(wt, "-list: invalid argument")
		return
	}

	var names []string
	for name := range ctrl.funcs {
		if !reg.MatchString(name) {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)

	ctrl.listFunctions(wt, names)
}

func (ctrl *Controller) listFunctions(wt io.Writer, names []string) {
	tw := tabwriter.NewWriter(wt, 0, 0, 2, ' ', 0)
	for _, name := range names {
		meta := ctrl.funcs[name]
		fmt.Fprintf(tw, "%s\t%s\t// %s\n", meta.name, meta.fn.Type(), meta.desc)
	}
	tw.Flush()

	if len(names) == 0 {
		printError(wt, "-list: function not found")
	}
}
