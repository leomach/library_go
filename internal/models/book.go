package models

import (
	"library_api/pkg/db"
)

type Book struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Author string  `json:"author"`
	Year   int     `json:"year"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
}

func (b *Book) CreateBook() error {
	sqlStatement := `
    INSERT INTO books (name, author, year, price, stock)
    VALUES (?,?,?,?,?)`

	stmt, err := db.DB.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(b.Name, b.Author, b.Year, b.Price, b.Stock)
	if err != nil {
		return err
	}
	return nil
}

func GetAllBooks() ([]Book, error) {
	rows, err := db.DB.Query("SELECT id, name, author, year, price, stock FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Name, &book.Author, &book.Year, &book.Price, &book.Stock)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func GetBookById(id string) (*Book, error) {
    var book Book

    // Ajuste a ordem do Scan para corresponder à ordem dos campos na consulta
    err := db.DB.QueryRow("SELECT price, stock FROM books WHERE id = ?", id).Scan(&book.Price, &book.Stock)
    if err != nil {
        return nil, err
    }

    return &book, nil
}


func DecrementBookStack(id string) error {
	sqlStatement := `UPDATE books SET stock = stock - 1 WHERE id = ?`
	stmt, err := db.DB.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
