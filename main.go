package main

import (
	"urlier/configs"
	"urlier/controllers"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// middleware logging
	e.Use(middleware.Logger())

	// connect to the database
	configs.ConnectDB()

	e.Static("/public", "public")

	e.GET("/", controllers.GetHome)
	e.GET("/about", controllers.GetAbout)
	e.GET("/tutorial", controllers.GetTutorial)
	e.GET("/empty", controllers.GetEmpty)
	e.GET("/trending", controllers.GetTrending)
	e.GET("/:key", controllers.GetKey)

	e.POST("/insert-key", controllers.PostInsertKey)

	e.Logger.Fatal(e.Start(":3000"))
}
