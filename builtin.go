package controller

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
)

func (this *Controller) registerBuiltInFuncs() {
	this.MustRegister(this.listFuncs, "_list", "return registered function list "+
		"(case insensitive, wildcard '*' supported)")
}

func (this *Controller) listFuncs(name string) (string, error) {
	metaList, err := filterFuncs(this.funcMap, name)
	if err != nil {
		return "", err
	}

	if len(metaList) == 0 {
		return "", errors.New("no results")
	}

	sort.Slice(metaList, func(i int, j int) bool {
		return metaList[i].name < metaList[j].name
	})
	sort.SliceStable(metaList, func(i int, j int) bool {
		return strings.ToLower(metaList[i].name) < strings.ToLower(metaList[j].name)
	})

	return formatFuncList(metaList), nil
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

func formatFuncList(metaList []*funcMeta) string {
	var buf bytes.Buffer

	tw := tabwriter.NewWriter(&buf, 0, 0, 4, ' ', 0)
	for _, meta := range metaList {
		fmt.Fprintf(tw, "\n%s\t%v\t// %s", meta.name, meta.fn.Type(), meta.desc)
	}
	tw.Flush()

	return buf.String()
}
