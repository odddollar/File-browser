package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	api := router.Group("/api")
	api.GET("*path", func(ctx *gin.Context) {
		files, err := os.ReadDir(ctx.Param("path")[1:])
		if err != nil {
			panic(err)
		}

		var response struct {
			Folders []string `json:"folders"`
			Files   []string `json:"files"`
		}

		for _, file := range files {
			if file.IsDir() {
				response.Folders = append(response.Folders, file.Name())
			} else {
				response.Files = append(response.Files, file.Name())
			}
		}

		ctx.JSON(200, response)
	})

	router.Run("localhost:8080")
}
