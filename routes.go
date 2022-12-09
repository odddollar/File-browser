package main

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Redirect from / home page to /app
func appRedirect(ctx *gin.Context) {
	ctx.Redirect(303, "/app")
}

// Return error 404 with appropriate template
func notFound(ctx *gin.Context) {
	ctx.HTML(404, "404.html", gin.H{
		"Message": "\"" + ctx.Request.Host + ctx.Request.URL.Path + "\" not found",
	})
}

// Return HTML template containing contents of given path
func viewDirectory(ctx *gin.Context) {
	// Get path from url and add to root path
	path := rootPath + ctx.Param("path")
	path = strings.ReplaceAll(path, "//", "/")

	// Read file path on server
	files, err := os.ReadDir(path)

	// Display 404 page if path not found
	if err != nil {
		notFound(ctx)
		return
	}

	// Create variable for storing directory information
	var response struct {
		URL     []string
		Path    string
		Folders []string
		Files   []string
	}

	// Add path and URL data to struct
	response.Path = path
	response.URL = deleteEmpty(strings.Split(strings.TrimPrefix(ctx.Request.URL.String(), "/app"), "/"))

	// Add file and folder information to struct
	for _, file := range files {
		if file.IsDir() {
			response.Folders = append(response.Folders, file.Name())
		} else {
			response.Files = append(response.Files, file.Name())
		}
	}

	// Send data to template
	ctx.HTML(200, "home.html", response)
}
