package commands

import (
	"fmt"
	"strings"
)

func Escape(s string) string {
	return strings.Replace(s, "'", "'\\''", -1)
}

func Echo(format string, a ...interface{}) string {
	s := fmt.Sprintf(format, a...)
	return fmt.Sprintf("echo \"%s\"\n", s)
}

func Command(format string, a ...interface{}) string {
	s := fmt.Sprintf(format, a...)
	return fmt.Sprintf("%s\n", s)
}

func Script(commands ...string) []byte {
	b := []byte("#!/bin/bash\n")
	for _, c := range commands {
		b = append(b, []byte(c)...)
	}
	return b
}
