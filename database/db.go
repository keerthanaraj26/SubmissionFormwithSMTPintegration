package database

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

var	DB  *sql.DB

func init() {
	var err error

	dsn := "root:test@123@tcp(127.0.0.1:3306)/student?"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
	}
}