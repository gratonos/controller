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
	msg := "type '-list [name]' to show registered functions " +
		"(case insensitive, wildcard '*' supported)"
	fmt.Fprintln(wt, colorize(msg, green))
}

func builtin(cmd string) bool {
	return strings.HasPrefix(cmd, "-")
}

func (ctrl *Controller) handleBuiltin(wt io.Writer, cmd string) {
	switch {
	case strings.HasPrefix(cmd, "-list"):
		ctrl.handleList(wt, strings.Fields(cmd)[1:])
	default:
		printError(wt, fmt.Sprintf("unsupported command '%s'", cmd))
	}
}

func (ctrl *Controller) handleList(wt io.Writer, args []string) {
	argc := len(args)
	if argc > 1 {
		printError(wt, fmt.Sprintf("-list: too many arguments, want 0 or 1, have %d", argc))
		return
	}

	arg, exp := "", ""
	if argc == 1 {
		arg = args[0]
		exp = `(?i)\b` + strings.Replace(arg, "*", ".*", -1) + `\b`
	}
	reg, err := regexp.Compile(exp)
	if err != nil {
		printError(wt, fmt.Sprintf("-list: invalid argument '%s'", arg))
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

	if len(names) == 0 {
		if argc == 0 {
			printError(wt, "-list: no registered functions")
		} else {
			printError(wt, fmt.Sprintf("-list: function '%s' not found", arg))
		}
	}
}

func (ctrl *Controller) listFunctions(wt io.Writer, names []string) {
	tw := tabwriter.NewWriter(wt, 0, 0, 4, ' ', 0)
	for _, name := range names {
		meta := ctrl.funcs[name]
		output := fmt.Sprintf("%s\t%v\t// %s", meta.name, meta.fn.Type(), meta.desc)
		fmt.Fprintln(tw, colorize(output, green))
	}
	tw.Flush()
}
