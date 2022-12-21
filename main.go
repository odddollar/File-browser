package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
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
	ginMode := parser.Flag("d", "dev", &argparse.Options{Default: false, Help: "Run Gin framework in debug/dev mode"})
	port := parser.Int("p", "port", &argparse.Options{Default: 8080, Help: "Port to host webserver on"})
	rP := parser.String("v", "path", &argparse.Options{Required: true, Help: "Root path to host. Must be absolute and exist"})

	// Run command line parser
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}

	// Set root path to CLI value
	rootPath = *rP
	rootPath = strings.ReplaceAll(rootPath, "\\", "/")

	// Check if root path is absolute to prevent weird path errors
	if !filepath.IsAbs(rootPath) {
		fmt.Println(parser.Usage("Path specified isn't absolute"))
		return
	}

	// Check that root path actually exists
	if !pathExists(rootPath) {
		fmt.Println(parser.Usage("Path specified doesn't exist"))
		return
	}

	// Create template
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"join":           strings.Join,
		"append":         templateAppend,
		"stripLastIndex": templateStripLastIndex,
		"isFile":         templateIsFile,
	}).ParseFS(fViews, "views/*.html", "views/*/*.html"))

	// Set release or debug mode
	if !(*ginMode) {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router and load HTML/static files
	router := gin.Default()
	router.SetHTMLTemplate(tmpl)
	router.StaticFS("/static", http.FS(subStatic(fStatic)))

	// Handle request to home page
	router.GET("/", appRedirect)

	// Handle path for directories and files
	router.GET("/app/*path", dirOrFile)

	// Router group for handling files
	file := router.Group("/file")
	{
		// Handle downloading files
		file.GET("/*path", downloadFile)

		// Handle postback for uploading/saving files
		file.POST("/*path", uploadFile)
	}

	// Router path for creating new items
	router.POST("/new/:type/*path", createNew)

	// Add route for 404
	router.NoRoute(notFound)

	// Run server
	router.Run(fmt.Sprintf("localhost:%d", *port))
}
