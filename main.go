package main

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const rootPath = "c:/users/sieea/documents"

//go:embed views
var fViews embed.FS

//go:embed static
var fStatic embed.FS

func main() {
	// create template
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"join":           strings.Join,
		"append":         templateAppend,
		"stripLastIndex": templateStripLastIndex,
	}).ParseFS(fViews, "views/*.html"))

	// create router and load HTML/static files
	router := gin.Default()
	router.SetHTMLTemplate(tmpl)
	router.StaticFS("/static", http.FS(subStatic(fStatic)))

	// handle request to home page, redirecting to appropriate URL
	router.GET("/", appRedirect)

	// handle path for viewing directories
	router.GET("/app/*path", viewDirectory)

	// add route for 404
	router.NoRoute(notFound)

	// run server
	router.Run("localhost:8080")
}
