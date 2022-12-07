package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const rootPath = "c:/users/sieea/documents"

func main() {
	// create router and load HTML files
	router := gin.Default()
	router.LoadHTMLGlob("views/*")

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
		response.URL = ctx.Request.URL.Path
		response.PreviousURL = strings.ReplaceAll(filepath.Dir(ctx.Request.URL.Path), "\\", "/")
		response.FileURL = strings.Replace(ctx.Request.URL.Path, "folder", "file", 1)
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

	// run server
	router.Run("localhost:8080")
}
