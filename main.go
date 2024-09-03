package main

import (
	"library_api/pkg/db"
	"library_api/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db.Connect()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/books", handlers.CreateBook)
	e.GET("/books", handlers.ListBooks)
	e.POST("/books/:id/buy", handlers.BuyBook)

	e.Logger.Fatal(e.Start(":1323"))
}
