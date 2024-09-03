package db


import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	var err error
	DB, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/library_db")
	if err != nil {
		panic(err)
	}
	
	fmt.Println("Database connected successfully")
}