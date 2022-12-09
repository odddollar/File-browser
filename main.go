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
	// Create template
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"join":           strings.Join,
		"append":         templateAppend,
		"stripLastIndex": templateStripLastIndex,
	}).ParseFS(fViews, "views/*.html"))

	// Create router and load HTML/static files
	router := gin.Default()
	router.SetHTMLTemplate(tmpl)
	router.StaticFS("/static", http.FS(subStatic(fStatic)))

	// Handle request to home page
	router.GET("/", appRedirect)

	// Handle path for viewing directories
	router.GET("/app/*path", viewDirectory)

	// Add route for 404
	router.NoRoute(notFound)

	// Run server
	router.Run("localhost:8080")
}
