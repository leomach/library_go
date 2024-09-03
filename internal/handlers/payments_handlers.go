package handlers

import (
	"database/sql"
	"fmt"
	"library_api/internal/models"
	"library_api/pkg/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func BuyBook(c echo.Context) error {
	id := c.Param("id")

	book, err := models.GetBookById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Error getting book")	
			return echo.NewHTTPError(http.StatusNotFound, "Book not found")
		}
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

	if book.Stock <= 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Book out of stock")
	}

	err = models.DecrementBookStack(id)
	if err != nil {
		fmt.Println("Error decrementing book stock")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	payment := new(models.Payment)
	if err := c.Bind(&payment); err != nil {
		fmt.Println("Error binding payment")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	payment.BookID = utils.Atoi(id)
	payment.PaymentValue = book.Price
	payment.PaymentDate = time.Now()

	err = payment.CreatePayment()
	if err != nil {
		fmt.Println("Error creating payment")
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	return c.JSON(http.StatusOK, payment)
}
