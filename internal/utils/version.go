package utils

import (
	"strconv"
	"strings"
)

func VersionToInt(version string) (major, minor, patch, build int) {
	if version != "" && version[0] == 'v' {
		items := strings.Split(version[1:], ".")
		if len(items) == 3 {
			major, _ = strconv.Atoi(items[0])
			minor, _ = strconv.Atoi(items[1])
			patch, _ = strconv.Atoi(items[2])
		}
		if len(items) == 4 {
			major, _ = strconv.Atoi(items[0])
			minor, _ = strconv.Atoi(items[1])
			patch, _ = strconv.Atoi(items[2])
			build, _ = strconv.Atoi(items[3])
		}
	}
	return
}
