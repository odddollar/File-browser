package main

import (
	"embed"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

// Template utility for checking if path is a file or directory.
// Used for determining what items/buttons to render in header
func templateIsFile(s []string) bool {
	// Join given path with root path
	path := rootPath + "/" + strings.Join(s, "/")

	// Get filesystem information
	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	// Return if a file
	if info.IsDir() {
		return false
	}
	return true
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

// Check if the given path is a subdirectory of the root path.
// Cleans path first
func isValidPath(path string) bool {
	// Normalise \ and /
	r := filepath.Clean(rootPath)

	// Clean given path
	p := filepath.Clean(path)

	// If the given path is in the cleaned root path, then the directory is valid
	return strings.Contains(p, r)
}

// Checks if a given path exists
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Only create file if it doesn't already exist.
// If is does exist, returns an error
func createFile(path string) error {
	// If file's path doesn't exist, create the file and return nil
	if !pathExists(path) {
		// If error occurs while writing to file, the path is likely malformed
		if err := os.WriteFile(path, []byte(""), 0755); err != nil {
			return err
		}

		// If the file path didn't exist and was created successfully, return nothing
		return nil
	}

	// Return error if file does exist
	return errors.New("Path already exists: " + path)
}
