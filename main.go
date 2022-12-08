package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const rootPath = "c:/users/sieea/documents"

func main() {
	// create router and load HTML/static files
	router := gin.Default()
	router.LoadHTMLGlob("views/*")
	router.Static("/static", "./static")

	// handle request to home page, redirecting to folder URL
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(303, "/folder/")
	})

	// handle path for viewing directories
	router.GET("/folder/*path", func(ctx *gin.Context) {
		// get path from url and add to root path
		path := rootPath + ctx.Param("path")

		// read file path on server
		files, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}

		// create variable for storing directory information
		var response struct {
			URL         string
			PreviousURL string
			FileURL     string
			Path        string
			Folders     []string
			Files       []string
		}

		// add relevant file and folder data to struct
		response.Path = path
		response.URL = strings.ReplaceAll(filepath.Clean(ctx.Request.URL.String()), "\\", "/")
		response.PreviousURL = strings.ReplaceAll(filepath.Dir(ctx.Request.URL.String()), "\\", "/")
		response.FileURL = strings.ReplaceAll(filepath.Clean(strings.Replace(ctx.Request.URL.String(), "folder", "file", 1)), "\\", "/")

		// add file and folder information to struct
		for _, file := range files {
			if file.IsDir() {
				response.Folders = append(response.Folders, file.Name())
			} else {
				response.Files = append(response.Files, file.Name())
			}
		}

		// send data to template
		ctx.HTML(200, "home.html", response)
	})

	// handle request to view individual file
	router.GET("/file/*path", func(ctx *gin.Context) {
		// get file system path
		path := rootPath + ctx.Param("path")

		// send file with appropriate name
		ctx.FileAttachment(path, filepath.Base(path))
	})

	router.POST("/file/*path", func(ctx *gin.Context) {
		// get equivalent folder url to redirect to
		folderURL := strings.ReplaceAll(filepath.Clean(strings.Replace(ctx.Request.URL.String(), "file", "folder", 1)), "\\", "/")

		// get file(s) from form data
		form, _ := ctx.MultipartForm()
		files := form.File["file"]

		// iterate through files in form
		for _, file := range files {
			// create path to save file to and save file
			path := strings.ReplaceAll(filepath.Clean(rootPath+ctx.Param("path")+"/"+file.Filename), "\\", "/")
			fmt.Println(file.Filename, path)
			ctx.SaveUploadedFile(file, path)
		}

		// redirect back to current page
		ctx.Redirect(303, folderURL)
	})

	// run server
	router.Run("localhost:8080")
}
