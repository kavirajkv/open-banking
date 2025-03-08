package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type accoutreq struct {
	Phone string `json:"phone"`
}

type accountres struct {
	Accno  string `json:"acc_no"`
	Ifsc   string `json:"ifsc"`
	Type   string `json:"acc_type"`
	Branch string `json:"branch"`
}

type transactionreq struct {
	Accno string `json:"acc_no"`
}

type transactionres struct {
	Accno        string      `json:"acc_no"`
	Ifsc         string      `json:"ifsc"`
	Transactions transaction `json:"transactions"`
}

type transaction struct {
	Transactionid int       `json:"transaction_id"`
	Type          string    `json:"transaction_type"`
	Mode          string    `json:"transaction_mode"`
	Amount        int       `json:"amount"`
	Time          time.Time `json:"time"`
}

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

func route() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/icici/get_account", get_account("icici")).Methods("POST", "OPTIONS")
	router.HandleFunc("/axis/get_account", get_account("axis")).Methods("POST", "OPTIONS")
	router.HandleFunc("/hdfc/get_account", get_account("hdfc")).Methods("POST", "OPTIONS")

	return router
}

func get_account(db string)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
	var accreq accoutreq

	err := json.NewDecoder(r.Body).Decode(&accreq)
	if err != nil {
		http.Error(w, "Enter string message", http.StatusBadRequest)
	}

	db := ConnectDB(db)
	defer db.Close()
	statement := "SELECT account_number,ifsc_code,account_type,branch_name FROM accounts WHERE user_id = (SELECT user_id FROM users WHERE phone_number = $1);"
	row := db.QueryRow(statement, accreq.Phone)
	var accres accountres
	errors := row.Scan(&accres.Accno, &accres.Ifsc, &accres.Type, &accres.Branch)
	if errors != nil {
		log.Fatalf("error while converting data -%v", errors)
	}

	json.NewEncoder(w).Encode(accres)

}
}

func main() {
	r := route()
	fmt.Println("Server running at port 8000")
	http.ListenAndServe(":8000", r)
}
