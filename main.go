package main

import (
	"html/template"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const rootPath = "c:/users/sieea/documents"

func main() {
	// create template
	tmpl := template.Must(template.New("main").Funcs(template.FuncMap{
		"join":           strings.Join,
		"append":         templateAppend,
		"stripLastIndex": templateStripLastIndex,
	}).ParseGlob("views/*.html"))

	// create router and load HTML/static files
	router := gin.Default()
	router.SetHTMLTemplate(tmpl)
	router.Static("/static", "./static")

	// handle request to home page, redirecting to appropriate URL
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(303, "/app")
	})

	// handle path for viewing directories
	router.GET("/app/*path", func(ctx *gin.Context) {
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
	})

	// // handle request to view individual file
	// router.GET("/file/*path", func(ctx *gin.Context) {
	// 	// get file system path
	// 	path := rootPath + ctx.Param("path")

	// 	// send file with appropriate name
	// 	ctx.FileAttachment(path, filepath.Base(path))
	// })

	// router.POST("/file/*path", func(ctx *gin.Context) {
	// 	// get equivalent folder url to redirect to
	// 	folderURL := strings.ReplaceAll(filepath.Clean(strings.Replace(ctx.Request.URL.String(), "file", "folder", 1)), "\\", "/")

	// 	// get file(s) from form data
	// 	form, _ := ctx.MultipartForm()
	// 	files := form.File["file"]

	// 	// iterate through files in form
	// 	for _, file := range files {
	// 		// create path to save file to and save file
	// 		path := strings.ReplaceAll(filepath.Clean(rootPath+ctx.Param("path")+"/"+file.Filename), "\\", "/")
	// 		ctx.SaveUploadedFile(file, path)
	// 	}

	// 	// redirect back to current page
	// 	ctx.Redirect(303, folderURL)
	// })

	// add route for 404
	router.NoRoute(notFound)

	// run server
	router.Run("localhost:8080")
}

func notFound(ctx *gin.Context) {
	ctx.HTML(404, "404.html", gin.H{
		"Message": "\"" + ctx.Request.Host + ctx.Request.URL.Path + "\" not found",
	})
}

func templateStripLastIndex(s []string) []string {
	return s[:len(s)-1]
}

func templateAppend(s []string, n string) []string {
	return append(s, n)
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
