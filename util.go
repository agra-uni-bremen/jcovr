package main

import (
	"os"
	"strings"
)

func getRelPath(n int) string {
	var ret string
	for i := 0; i < n; i++ {
		ret += "../"
	}

	if ret == "" {
		return "./"
	} else {
		return ret
	}
}

func relIndex(path string) string {
	elems := strings.SplitN(path, string(os.PathSeparator), -1)
	return getRelPath(len(elems) - 1)
}
