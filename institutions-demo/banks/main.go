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

// account discovery usiing phone type
type accountreq struct {
	Phone string `json:"phone"`
}

// if accound found return type for response
type accountres struct {
	Accno  string `json:"acc_no"`
	Ifsc   string `json:"ifsc"`
	Type   string `json:"acc_type"`
	Branch string `json:"branch"`
}

// transactions request type
type transactionreq struct {
	Accno string `json:"acc_no"`
	Bank  string `json:"bank"`
}

// return type for transactions data request
type transactionres struct {
	Accno        string        `json:"acc_no"`
	Transactions []transaction `json:"transactions"`
}

type transaction struct {
	Transactionid int       `json:"transaction_id"`
	Type          string    `json:"transaction_type"`
	Mode          string    `json:"transaction_mode"`
	Amount        float64   `json:"amount"`
	Time          time.Time `json:"time"`
}

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

// routes
func route() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/icici/get_account", get_account("icici")).Methods("POST", "OPTIONS")
	router.HandleFunc("/axis/get_account", get_account("axis")).Methods("POST", "OPTIONS")
	router.HandleFunc("/hdfc/get_account", get_account("hdfc")).Methods("POST", "OPTIONS")
	router.HandleFunc("/get_transactions", get_transactions).Methods("POST", "OPTIONS")

	return router
}

// this function will search for account based on user phone number if found return details
func get_account(db string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var accreq accountreq

		err := json.NewDecoder(r.Body).Decode(&accreq)
		if err != nil {
			http.Error(w, "Enter valid phone number", http.StatusBadRequest)
		}

		db := ConnectDB(db)
		defer db.Close()
		statement := "SELECT account_number,ifsc_code,account_type,branch_name FROM accounts WHERE user_id = (SELECT user_id FROM users WHERE phone_number = $1);"
		row := db.QueryRow(statement, accreq.Phone)
		var accres accountres
		errors := row.Scan(&accres.Accno, &accres.Ifsc, &accres.Type, &accres.Branch)
		if errors != nil {
			http.Error(w, "account not found", http.StatusNotFound)

		}

		json.NewEncoder(w).Encode(accres)

	}
}

// this function will retern all the transaction based on acc number
func get_transactions(w http.ResponseWriter, r *http.Request) {
	var transactionrequest transactionreq

	err := json.NewDecoder(r.Body).Decode(&transactionrequest)
	if err != nil {
		http.Error(w, "Enter valid account number", http.StatusBadRequest)
	}

	db := ConnectDB(transactionrequest.Bank)
	defer db.Close()
	statement := "select transaction_id,transaction_type,mode,amount,transaction_timestamp from transactions where account_id=(select account_id from accounts where account_number=$1);"
	rows, _ := db.Query(statement, transactionrequest.Accno)
	var response transactionres

	var transactions []transaction
	for rows.Next() {
		var trans transaction
		errors := rows.Scan(&trans.Transactionid, &trans.Type, &trans.Mode, &trans.Amount, &trans.Time)
		if errors != nil {
			http.Error(w, "account not found", http.StatusNotFound)

		}
		transactions = append(transactions, trans)

	}
	response.Accno = transactionrequest.Accno
	response.Transactions = transactions

	json.NewEncoder(w).Encode(response)

}

func main() {
	r := route()
	fmt.Println("Server running at port 8000")
	http.ListenAndServe(":8000", r)
}
