package main

import (
	"html/template"
	"strings"

	"github.com/gin-gonic/gin"
)

const rootPath = "c:/users/sieea/documents"

func main() {
	// create template
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"join":           strings.Join,
		"append":         templateAppend,
		"stripLastIndex": templateStripLastIndex,
	}).ParseGlob("views/*.html"))

	// create router and load HTML/static files
	router := gin.Default()
	router.SetHTMLTemplate(tmpl)
	router.Static("/static", "./static")

	// handle request to home page, redirecting to appropriate URL
	router.GET("/", appRedirect)

	// handle path for viewing directories
	router.GET("/app/*path", viewDirectory)

	// add route for 404
	router.NoRoute(notFound)

	// run server
	router.Run("localhost:8080")
}
