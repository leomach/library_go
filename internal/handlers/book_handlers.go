package handlers

import (
    "github.com/labstack/echo/v4"
	"library_api/internal/models"
	"net/http"
)

func CreateBook(c echo.Context) error {
	b := new(models.Book)
	if err := c.Bind(b); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}

	err := b.CreateBook()
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError)
    }

	return c.JSON(http.StatusCreated, b)
}

func ListBooks(c echo.Context) error {
	books, err := models.GetAllBooks()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	return c.JSON(http.StatusOK, books)
}