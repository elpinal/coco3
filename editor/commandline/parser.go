package commandline

import "strings"

func Parse(s string) []string {
	return strings.Split(s, " ")
}
