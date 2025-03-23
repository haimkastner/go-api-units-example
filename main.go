package main

import (
	"embed"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	routes "github.com/haimkastner/go-api-units-example/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed openapi/openapi.json
var swaggerSpec embed.FS

// Handler to serve the Swagger spec file
func serveSwaggerSpec(c *gin.Context) {
	// Read the embedded swagger.json file
	specData, err := swaggerSpec.ReadFile("openapi/openapi.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read openapi spec",
		})
		return
	}

	// Set the content type to application/json
	c.Header("Content-Type", "application/json")
	// Write the spec data directly to the response
	c.Writer.Write(specData)
}

func main() {

	// Create a default gin router
	router := gin.Default()

	// Add CORS middleware with permissive settings
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"https://units-docs.gleece.dev"}, // Allow all origins
		AllowMethods:  []string{"*"},                             // Allow all methods
		AllowHeaders:  []string{"*"},                             // Allow all headers
		ExposeHeaders: []string{"Content-Length"},
	}))

	// Register custom validation rules
	routes.RegisterRoutes(router)

	// # End Gleece integration part

	// Serve the Swagger spec file at /openapi/openapi.json
	router.GET("/openapi/openapi.json", serveSwaggerSpec)

	// Serve the Swagger UI at /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/openapi/openapi.json")))

	// Get port from environment variable or use 8080 as default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	router.Run(":" + port)
}
