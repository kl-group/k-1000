package main

import (
	"regexp"
	"strings"
)

func helpersGetNameFromDN(cn string) string {
	split := strings.Split(cn, ",")
	if len(split) == 0 {
		return cn
	}
	var re = regexp.MustCompile("^CN=(.*)$")
	s := re.ReplaceAllString(split[0], `$1`)
	return s
}

func reloadAsterisk() {
	if needReloadAsterisk == false {
		return
	}
	err := GraphQLReloadAsterisk()
	if err != nil {
		logme.Warning(err)
	}
}
