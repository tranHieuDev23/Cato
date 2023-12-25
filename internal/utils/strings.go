package utils

import "strings"

func TrimSpaceRight(s string) string {
	return strings.TrimRight(s, "\t\n\r\v\f ")
}
