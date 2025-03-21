package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)


// db connection based on database name
func ConnectDB(dbname string) *sql.DB {
	db_pass := os.Getenv("PG_PASSWORD")
	dsn := fmt.Sprintf("postgres://postgres:%v@localhost/%v?sslmode=disable", db_pass, dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	fmt.Println("Database connected..")
	return db

}