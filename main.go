package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Book struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Author string  `json:"author"`
	Year   int     `json:"year"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
}

type Payment struct {
	ID            int       `json:"id"`
	BookID        int       `json:"book_id"`
	PaymentValue  float64   `json:"payment_value"`
	PaymentMethod string    `json:"payment_method"`
	PaymentDate   time.Time `json:"payment_date"`
}

var (
	db *sql.DB
)

func createBook(c echo.Context) error {
	b := new(Book)
	if err := c.Bind(b); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}

	sqlStatement := `
    INSERT INTO books (name, author, year, price, stock)
    VALUES (?,?,?,?,?)`

	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway)
	}
	defer stmt.Close()

	_, err = stmt.Exec(b.Name, b.Author, b.Year, b.Price, b.Stock)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway)
	}
	
	return c.JSON(http.StatusCreated, b)
}

func listBooks(c echo.Context) error {
	rows, err := db.Query("SELECT id, name, author, year, price, stock FROM books")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Name, &book.Author, &book.Year, &book.Price, &book.Stock)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, books)
}

func buyBook(c echo.Context) error {
	id := c.Param("id")
	var stock int
	var price float64
	err := db.QueryRow("SELECT stock, price FROM books WHERE id = ?", id).Scan(&stock, &price)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Book not found")
		}
		return echo.NewHTTPError(http.StatusBadGateway)
	}

	if stock <= 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Book out of stock")
	}

	sqlStatement := `UPDATE books SET stock = stock - 1 WHERE id = ?`
	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway)
	}

	payment := new(Payment)
	if err := c.Bind(payment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payment input"})
	}
	payment.BookID = atoi(id)
	payment.PaymentValue = price
	payment.PaymentDate = time.Now()

	sqlPayment := `
    INSERT INTO payments (book_id, payment_value, payment_method, payment_date)
    VALUES (?, ?, ?, NOW())`

	stmtPay, err := db.Prepare(sqlPayment)
	if err != nil {
		return err
	}
	defer stmtPay.Close()
	_, err = stmtPay.Exec(payment.BookID, payment.PaymentValue, payment.PaymentMethod)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	
	return c.JSON(http.StatusOK, payment)
}

func atoi(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/library_db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Database connected successfully")

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/books", createBook)
	e.GET("/books", listBooks)
	e.POST("/books/:id/buy", buyBook)

	e.Logger.Fatal(e.Start(":1323"))
}
