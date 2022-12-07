package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*")

	api := router.Group("/api")
	api.GET("*path", func(ctx *gin.Context) {
		path := ctx.Param("path")[1:]

		files, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}

		var response struct {
			Path    string   `json:"path"`
			Folders []string `json:"folders"`
			Files   []string `json:"files"`
		}

		response.Path = path
		for _, file := range files {
			if file.IsDir() {
				response.Folders = append(response.Folders, file.Name())
			} else {
				response.Files = append(response.Files, file.Name())
			}
		}

		ctx.JSON(200, response)
	})

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "home.html", gin.H{})
	})

	router.Run("localhost:8080")
}
