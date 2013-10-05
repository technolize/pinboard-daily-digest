package main

import (
	"os"
	"strings"
)

func Environ() map[string]string {
	env := make(map[string]string)

	for _, item := range os.Environ() {
		parts := strings.Split(item, "=")
		key := parts[0]
		val := strings.Join(parts[1:], "")
		env[key] = val
	}

	return env
}
