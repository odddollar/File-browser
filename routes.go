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
		notFound(ctx)
		return
	}

	// Run handler function for directory or path
	if info.IsDir() {
		viewDirectory(ctx, path)
	} else {
		viewFile(ctx, path)
	}
}

// View file with text editor
func viewFile(ctx *gin.Context, path string) {
	// Read file to string
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	content := string(file)

	// Get split URL of file
	// Used by header buttons to determine where to send form data
	URL := deleteEmpty(strings.Split(ctx.Param("path"), "/"))

	// Send data to template
	ctx.HTML(200, "edit.html", gin.H{"Content": content, "URL": URL, "Path": path})
}

// Return HTML template containing contents of given path
func viewDirectory(ctx *gin.Context, path string) {
	// Read file path on server
	files, err := os.ReadDir(path)

	// Display 404 page if path not found
	if err != nil {
		panic(err)
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
	response.URL = deleteEmpty(strings.Split(ctx.Param("path"), "/"))

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

// Download file from given path
func downloadFile(ctx *gin.Context) {
	// Create path to file
	path := rootPath + ctx.Param("path")
	path = strings.ReplaceAll(path, "//", "/")

	// Send file as attachment
	ctx.FileAttachment(path, filepath.Base(path))
}

// Upload file to server
func uploadFile(ctx *gin.Context) {
	// Run based on content type header
	// Will be json if saving file, multipart form if uploading file
	if ctx.ContentType() == "application/json" {
		// Bind json data to variable
		var jsonData struct {
			Content string `json:"fileContent"`
		}
		ctx.BindJSON(&jsonData)

		// Get contents and path of file
		path := rootPath + ctx.Param("path")

		// Write new contents to file
		err := os.WriteFile(path, []byte(jsonData.Content), 0755)
		if err != nil {
			panic(err)
		}
	} else {
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
	}

	// Return successful status
	ctx.Status(200)
}

// Create new folder on server
func createNewFolder(ctx *gin.Context) {
	// Bind json data to variable
	var jsonData struct {
		Name string `json:"name"`
	}
	ctx.BindJSON(&jsonData)

	// Create full new folder path
	path := rootPath + ctx.Param("path") + "/" + jsonData.Name
	path = strings.ReplaceAll(path, "//", "/")

	// Check that path is valid and doesn't escape root path
	if isValidPath(path) {
		// Make path with set permissions
		err := os.Mkdir(path, 0755)
		if err != nil {
			// Server was not able create this directory.
			// This is mainly used if a path is valid, but malformed
			// (i.e. "C:/Windows/C:/users" where "C:/Windows" is the root path),
			// which helps prevent creating folders outside the root path.
			// Also used to prevent creating folders with names that already exist
			ctx.Status(403)
			return
		}
	} else {
		// Path isn't valid and user was likely trying to escape the root path
		ctx.Status(403)
		return
	}

	// Return successful status
	ctx.Status(200)
}

// Create new file on server
func createNewFile(ctx *gin.Context) {
	// Bind json data to variable
	var jsonData struct {
		Name string `json:"name"`
	}
	ctx.BindJSON(&jsonData)

	// Create full new file path
	path := rootPath + ctx.Param("path") + "/" + jsonData.Name
	path = strings.ReplaceAll(path, "//", "/")

	// Check that path is valid and doesn't escape root path
	if isValidPath(path) {
		// Make file with set permissions
		err := createFile(path)
		if err != nil {
			// Server was not able create this file.
			// This is mainly used if a path is valid, but malformed
			// (i.e. "C:/Windows/C:/users/hello.txt" where "C:/Windows" is the root path),
			// which helps prevent creating files outside the root path.
			// Also used to prevent creating files with names that already exist
			ctx.Status(403)
			return
		}
	} else {
		// Path isn't valid and user was likely trying to escape the root path
		ctx.Status(403)
		return
	}

	// Return successful status
	ctx.Status(200)
}
