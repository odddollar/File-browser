package main

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func appRedirect(ctx *gin.Context) {
	ctx.Redirect(303, "/app")
}

func notFound(ctx *gin.Context) {
	ctx.HTML(404, "404.html", gin.H{
		"Message": "\"" + ctx.Request.Host + ctx.Request.URL.Path + "\" not found",
	})
}

func viewDirectory(ctx *gin.Context) {
	// get path from url and add to root path
	path := rootPath + ctx.Param("path")

	// read file path on server
	files, err := os.ReadDir(path)
	if err != nil {
		notFound(ctx)
		return
	}

	// create variable for storing directory information
	var response struct {
		URL     []string
		Path    string
		Folders []string
		Files   []string
	}

	// add path and URL data to struct
	response.Path = path
	response.URL = deleteEmpty(strings.Split(strings.TrimPrefix(ctx.Request.URL.String(), "/app"), "/"))

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
}
