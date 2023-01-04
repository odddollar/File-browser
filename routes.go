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
	ctx.HTML(404, "error.html", gin.H{
		"Error":   404,
		"Message": "\"" + ctx.Request.Host + ctx.Request.URL.Path + "\" not found",
	})
}

// Return error 403 and template stating not permitted to access that file/directory
func notPermitted(ctx *gin.Context) {
	// Get split URL of file
	URL := deleteEmpty(strings.Split(ctx.Param("path"), "/"))

	ctx.HTML(403, "error.html", gin.H{
		"Error":   403,
		"Message": "\"" + ctx.Request.Host + ctx.Request.URL.Path + "\" is not accessible. Permission is likely denied",
		"URL":     URL,
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

	// Display 403 page if not able to read
	if err != nil {
		notPermitted(ctx)
		return
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

	// Display 403 page if not able to read
	if err != nil {
		notPermitted(ctx)
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

// Create new file or folder on server
func createNew(ctx *gin.Context) {
	// Bind json data to variable
	var jsonData struct {
		Name string `json:"name"`
	}
	ctx.BindJSON(&jsonData)

	// Create full new path
	path := rootPath + ctx.Param("path") + "/" + jsonData.Name
	path = strings.ReplaceAll(path, "//", "/")

	// Check that path is valid and doesn't escape root path
	if isValidPath(path) {
		var err error

		// Check whether file or folder needs to be created
		if ctx.Param("type") == "folder" {
			// Create folder
			err = os.Mkdir(path, 0755)
		} else if ctx.Param("type") == "file" {
			// Create file
			err = createFile(path)
		}

		if err != nil {
			// Server was not able create this file or folder.
			// This is mainly used if a path is valid, but malformed
			// (i.e. "C:/Windows/C:/users/hello.txt" where "C:/Windows" is the root path),
			// which helps prevent creating things outside the root path.
			// Also used to prevent creating things with names that already exist
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
