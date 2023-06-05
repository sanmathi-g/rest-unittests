package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	controllers "github.com/sana/rest/controller"
	"github.com/sana/rest/database"
)

func main() {
	fmt.Println("Starting application ...")
	database.DatabaseConnection()

	r := gin.Default()
	r.GET("/movies/:id", controllers.ReadMovie)
	r.GET("/movies", controllers.ReadMovies)
	r.POST("/movies", controllers.CreateMovie)
	r.PUT("/movies/:id", controllers.UpdateMovie)
	r.DELETE("/movies/:id", controllers.DeleteMovie)
	r.Run(":8080")
}
