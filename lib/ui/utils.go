package ui

import (
	"regexp"
	"strings"

	"github.com/enescakir/emoji"
	"github.com/gizak/termui/v3/widgets"
	lib "github.com/vinodsr/shell-butler/lib/types"
)

func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func updateContextList(contextList []string, contextFilter string) []string {

	filteredList := filter(contextList, func(str string) bool {
		// return strings.Contains(str, sel)

		match, _ := regexp.MatchString(contextFilter, str)
		return match
	})
	return filteredList
}

func displayContextatLevel(contextList []string, level int) []string {
	result := []string{}
	contextMap := make(map[string]bool)
	for _, s := range contextList {
		splits := strings.Split(s, ":")
		if len(splits) >= level {
			if contextMap[splits[level-1]] == false {
				result = append(result, splits[level-1])
				contextMap[splits[level-1]] = true
			}

		}
	}
	return result
}

func debug(str string, d *widgets.Paragraph) {
	d.Text += str + "\n"
}

func formatCommandString(selectedContext []string, commandStr string) string {
	formattedCommandStr := " "
	for _, s := range selectedContext {
		formattedCommandStr += "[" + s + "] "
	}
	return formattedCommandStr + emoji.RightArrow.String() + " " + commandStr
}

func RemoveIndex(s []lib.Command, index int) []lib.Command {
	return append(s[:index], s[index+1:]...)
}
