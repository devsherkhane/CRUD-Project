package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	// Database credentials from your original db.go
	dsn := "root:9595520628@@tcp(127.0.0.1:3306)/student_db"
	DB, err = sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("DB open error:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	fmt.Println("✅ MySQL Connected")
}