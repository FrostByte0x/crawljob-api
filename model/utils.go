package model

import (
	"strconv"
	"strings"
)

func booltoString(b bool) string {
	return strings.ToUpper(strconv.FormatBool(b))
}
