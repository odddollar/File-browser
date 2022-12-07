package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

const rootPath = "c:/users/sieea/documents"

func main() {
	// create router and load HTML files
	router := gin.Default()
	router.LoadHTMLGlob("views/*")

	// handle path for viewing directory
	router.GET("/*path", func(ctx *gin.Context) {
		// get path from url and add to root path
		path := rootPath + ctx.Param("path")

		// read file path on server
		files, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}

		// create variable for storing directory information
		var response struct {
			Path    string
			Folders []string
			Files   []string
		}

		// add relevant file and folder data to struct
		response.Path = path
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

	// run server
	router.Run("localhost:8080")
}
