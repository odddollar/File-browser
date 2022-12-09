package main

import (
	"embed"
	"io/fs"
)

func templateStripLastIndex(s []string) []string {
	return s[:len(s)-1]
}

func templateAppend(s []string, n string) []string {
	return append(s, n)
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func subStatic(f embed.FS) fs.FS {
	t, _ := fs.Sub(f, "static")
	return t
}
