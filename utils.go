package main

import (
	"embed"
	"io/fs"
)

// Template utility for removing final index in array of strings
func templateStripLastIndex(s []string) []string {
	return s[:len(s)-1]
}

// Template utility for appending string to array of strings
// (for some reason it doesn't like the regular "append" function)
func templateAppend(s []string, n string) []string {
	return append(s, n)
}

// Remove indexes in array of strings that contain an empty string
func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// Convert /static/static/* path for embedded files to /static/*
func subStatic(f embed.FS) fs.FS {
	t, _ := fs.Sub(f, "static")
	return t
}
