package controller

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

func (this *Controller) handleBuiltin(writer io.Writer, cmd string) {
	switch {
	case strings.HasPrefix(cmd, "-list"):
		this.handleList(writer, strings.Fields(cmd)[1:])
	default:
		printError(writer, "unsupported command '%s'", cmd)
	}
}

func (this *Controller) handleList(writer io.Writer, args []string) {
	metaList, err := this.filterFuncs(args)
	if err != nil {
		printError(writer, "%v", err)
		return
	}

	sort.Slice(metaList, func(i int, j int) bool {
		return metaList[i].name < metaList[j].name
	})

	printFuncList(writer, metaList)
}

func (this *Controller) filterFuncs(args []string) ([]*funcMeta, error) {
	argc := len(args)
	if argc > 1 {
		return nil, fmt.Errorf("-list: too many arguments, want 0 or 1, have %d", argc)
	}

	arg, exp := "", ""
	if argc == 1 {
		arg = args[0]
		exp = `(?i)\b` + strings.Replace(arg, "*", ".*", -1) + `\b`
	}
	reg, err := regexp.Compile(exp)
	if err != nil {
		return nil, fmt.Errorf("-list: invalid argument '%s'", arg)
	}

	var metaList []*funcMeta
	for name, meta := range this.funcs {
		if !reg.MatchString(name) {
			continue
		}
		metaList = append(metaList, meta)
	}

	if len(metaList) == 0 {
		if argc == 0 {
			return nil, errors.New("-list: no registered functions")
		} else {
			return nil, fmt.Errorf("-list: function '%s' not found", arg)
		}
	}

	return metaList, nil
}

func builtin(cmd string) bool {
	return strings.HasPrefix(cmd, "-")
}
