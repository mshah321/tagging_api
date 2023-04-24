package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	key := os.Getenv("KEY")
	google_tags := read_tags_from_file("google_tags.txt")

	get_tags_controller := NearByTagsController{google_tags, key}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Temporarily alllow all origins, methods, and headers. Restrict this.
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(config))

	// Router
	r.GET("/get_tags", get_tags_controller.get_nearby_tags)
	r.GET("/ping", func(ctx *gin.Context) { ctx.Status(200) })

	r.Run()
}
