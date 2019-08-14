package controller

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

func (this *Controller) handleBuiltinCmd(writer io.Writer, input string) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		panic("checking input failure")
	}

	cmd, args := fields[0], fields[1:]
	switch cmd {
	case "-list":
		this.handleCmdList(writer, args)
	default:
		printError(writer, "unsupported command '%s'", cmd)
	}
}

func (this *Controller) handleCmdList(writer io.Writer, args []string) {
	argc := len(args)
	if argc > 1 {
		printError(writer, "-list: too many arguments, want 0 or 1, have %d", argc)
		return
	}

	var arg string
	if argc == 1 {
		arg = args[0]
	}

	metaList, err := filterFuncs(this.funcMap, arg)
	if err != nil {
		printError(writer, "-list: %v", err)
		return
	}

	if len(metaList) == 0 {
		printError(writer, "-list: no results")
		return
	}

	sort.Slice(metaList, func(i int, j int) bool {
		return metaList[i].name < metaList[j].name
	})

	printFuncList(writer, metaList)
}

func isBuiltinCmd(input string) bool {
	return strings.HasPrefix(input, "-")
}

func filterFuncs(funcMap map[string]*funcMeta, arg string) ([]*funcMeta, error) {
	var exp string
	if arg != "" {
		exp = `(?i)\b` + strings.Replace(arg, "*", ".*", -1) + `\b`
	}

	reg, err := regexp.Compile(exp)
	if err != nil {
		return nil, fmt.Errorf("invalid argument '%s'", arg)
	}

	var metaList []*funcMeta
	for name, meta := range funcMap {
		if reg.MatchString(name) {
			metaList = append(metaList, meta)
		}
	}

	return metaList, nil
}
