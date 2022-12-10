package main

import (
	"os"
	"path/filepath"
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

// Check if request if for directory or file, handing context to relevant function
func dirOrFile(ctx *gin.Context) {
	// Get path from url and add to root path
	path := rootPath + ctx.Param("path")
	path = strings.ReplaceAll(path, "//", "/")

	// Get path information
	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	// Run handler function for directory or path
	if info.IsDir() {
		viewDirectory(ctx, path)
	} else {
		viewFile(ctx, path)
	}
}

// Download file for viewing
func viewFile(ctx *gin.Context, path string) {
	ctx.FileAttachment(path, filepath.Base(path))
}

// Return HTML template containing contents of given path
func viewDirectory(ctx *gin.Context, path string) {
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

// Upload file to server and redirect to original page
func uploadFile(ctx *gin.Context) {
	// Get list of files in form data
	form, _ := ctx.MultipartForm()
	files := form.File["file"]

	// Process each file
	for _, file := range files {
		// Create save path location
		path := rootPath + ctx.Param("path") + "/" + file.Filename
		path = strings.ReplaceAll(path, "//", "/")

		// Save current file
		ctx.SaveUploadedFile(file, path)
	}

	// Redirect back to original page
	ctx.Redirect(303, strings.Replace(ctx.Request.URL.String(), "file", "app", 1))
}

// Create new folder on server and redirect to original page
func createNewFolder(ctx *gin.Context) {
	// Create full new folder path
	path := rootPath + ctx.Param("path") + "/" + ctx.PostForm("new-folder-name")
	path = strings.ReplaceAll(path, "//", "/")

	// Make path with set permissions
	err := os.Mkdir(path, 0755)
	if err != nil {
		panic(err)
	}

	// Redirect back to original page
	ctx.Redirect(303, "/app"+ctx.Param("path"))
}

// Create new file on server and redirect to original page
func createNewFile(ctx *gin.Context) {
	// Create full new file path
	path := rootPath + ctx.Param("path") + "/" + ctx.PostForm("new-file-name")
	path = strings.ReplaceAll(path, "//", "/")

	// Make file with set permissions
	err := os.WriteFile(path, []byte(""), 0755)
	if err != nil {
		panic(err)
	}

	// Redirect back to original page
	ctx.Redirect(303, "/app"+ctx.Param("path"))
}
