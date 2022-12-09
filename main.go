package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/gin-gonic/gin"
)

//go:embed views
var fViews embed.FS

//go:embed static
var fStatic embed.FS

// Global variable to keep track of root path
var rootPath string

func main() {
	// Setup command line arguments
	parser := argparse.NewParser("Network File Browser", "View file system contents over a network")
	port := parser.Int("p", "port", &argparse.Options{Default: 8080, Help: "Port to host webserver on"})
	rP := parser.String("v", "path", &argparse.Options{Required: true, Help: "Root path to host"})

	// Run command line parser
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}

	// Set root path to CLI value
	rootPath = *rP

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

	// Handle path for directories and files
	router.GET("/app/*path", dirOrFile)

	// Add route for 404
	router.NoRoute(notFound)

	// Run server
	router.Run(fmt.Sprintf("localhost:%d", *port))
}
