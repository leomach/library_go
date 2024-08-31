package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	ID     int
	Name   string
	Author string
	Year   int
	Price  float64
	Stock  int
}

type Payment struct {
	ID            int
	BookID        int
	PaymentValue  float64
	PaymentMethod string
}

var (
	db *sql.DB
)

func main() {
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/library_db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Database connected successfully")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nEscolha uma opção:")
		fmt.Println("1. Inserir um novo livro")
		fmt.Println("2. Listar todos os livros")
		fmt.Println("3. Comprar um livro")
		fmt.Println("0. Sair")

		var choice int
		fmt.Print("Escolha: ")
		fmt.Scan(&choice)
		scanner.Scan()

		switch choice {
		case 1:
			var book Book
			fmt.Println("Inserir um novo livro")
			
			fmt.Print("Nome: ")
			book.Name = readString(scanner)
			fmt.Print("Autor: ")
			book.Author = readString(scanner)
			fmt.Print("Ano: ")
			book.Year = readInt(scanner)
			fmt.Print("Preço: ")
			book.Price = readFloat(scanner)
			fmt.Print("Estoque: ")
			book.Stock = readInt(scanner)

			err = insertBook(book)
			if err != nil {
				log.Println("Erro ao inserir o livro:", err)
			}

		case 2:
			fmt.Println("Listar todos os livros")
			books, err := getAllBooks()
			if err != nil {
				log.Println("Erro ao listar os livros:", err)
				continue
			}

			for _, book := range books {
				fmt.Printf("ID: %d, Nome: %s, Autor: %s, Ano: %d, Preço: %.2f, Estoque: %d\n",
					book.ID, book.Name, book.Author, book.Year, book.Price, book.Stock)
			}

		case 3:
			var bookID int
			var paymentMethod string
			fmt.Println("Comprar um livro")
			fmt.Print("ID do livro: ")
			fmt.Scan(&bookID)
			
			fmt.Println("Métodos de pagamento disponíveis:")
			fmt.Println("1. Crédito")
			fmt.Println("2. Débito")
			fmt.Println("3. Boleto")
			fmt.Println("4. Pix")
			fmt.Println("5. Dinheiro")
			fmt.Print("Escolha o método de pagamento (1-5): ")
			var methodChoice int
			fmt.Scan(&methodChoice)

			switch methodChoice {
			case 1:
				paymentMethod = "Crédito"
			case 2:
				paymentMethod = "Débito"
			case 3:
				paymentMethod = "Boleto"
			case 4:
				paymentMethod = "Pix"
			case 5:
				paymentMethod = "Dinheiro"
			default:
				fmt.Println("Método de pagamento inválido.")
				continue
			}

			err := buyBook(bookID, paymentMethod)
			if err != nil {
				log.Println("Erro ao comprar o livro:", err)
			}

		case 0:
			fmt.Println("Saindo...")
			os.Exit(0)

		default:
			fmt.Println("Opção inválida. Tente novamente.")
		}
	}
}

func readString(scanner *bufio.Scanner) string {
	scanner.Scan()
	return scanner.Text()
}

func readInt(scanner *bufio.Scanner) int {
	var input string
	scanner.Scan()
	input = scanner.Text()
	value, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Erro ao ler o número:", err)
		return 0
	}
	return value
}

func readFloat(scanner *bufio.Scanner) float64 {
	var input string
	scanner.Scan()
	input = scanner.Text()
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Println("Erro ao ler o número:", err)
		return 0
	}
	return value
}

func insertBook(book Book) error {
	sqlStatement := `
    INSERT INTO books (name, author, year, price, stock)
    VALUES (?,?,?,?,?)`

	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(book.Name, book.Author, book.Year, book.Price, book.Stock)
	if err != nil {
		return err
	}
	fmt.Println("New book inserted successfully")
	return nil
}

func getAllBooks() ([]Book, error) {
	rows, err := db.Query("SELECT id, name, author, year, price, stock FROM books")
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

func buyBook(bookID int, paymentMethod string) error {
	var stock int
	var price float64
	err := db.QueryRow("SELECT stock, price FROM books WHERE id = ?", bookID).Scan(&stock, &price)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Book not found")
		}
		return err
	}

	if stock <= 0 {
		return fmt.Errorf("Book out of stock")
	}

	sqlStatement := `UPDATE books SET stock = stock - 1 WHERE id = ?`
	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(bookID)
	if err!= nil {
        return err
    }
	fmt.Printf("Book %d purchased successfully. Remaining stock: %d\n", bookID, stock-1)

	sqlPayment := `
    INSERT INTO payments (book_id, payment_value, payment_method, payment_date)
    VALUES (?, ?, ?, NOW())`

	stmtPay, err := db.Prepare(sqlPayment)
	if err != nil {
		return err
	}
	defer stmtPay.Close()
	_, err = stmtPay.Exec(bookID, price, paymentMethod)
	if err!= nil {
        return err
    }
	fmt.Printf("Payment for book %d recorded successfully\n", bookID)
	return nil
}
