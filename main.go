package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// create router and load HTML files
	router := gin.Default()
	router.LoadHTMLGlob("views/*")

	// create router group for handling api
	api := router.Group("/api")
	api.GET("*path", func(ctx *gin.Context) {
		// get path from url
		path := ctx.Param("path")[1:]

		// read file path on server
		files, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}

		// create variable for converting to json
		var response struct {
			Path    string   `json:"path"`
			Folders []string `json:"folders"`
			Files   []string `json:"files"`
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

		// return struct as json
		ctx.JSON(200, response)
	})

	// host main html page
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "home.html", gin.H{})
	})

	// run server
	router.Run("localhost:8080")
}
