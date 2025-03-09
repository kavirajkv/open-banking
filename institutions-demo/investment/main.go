package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type investreq struct {
	Phone string `json:"phone"`
}

type investres struct {
	Demat_no string `json:"demat_no"`
	Broker   string `json:"broker"`
}

// response type to get mutual fund account
type mfresponse struct {
	Folio_no string `json:"folio_no"`
	Broker   string `json:"broker"`
}

// stocks details request type
type stocksreq struct {
	Demat_no string `json:"demat_no"`
}

// stock details response type
type stocksres struct {
	Demat_no string         `json:"demat_no"`
	Stocks   []stockdetails `json:"stocks"`
}

type stockdetails struct {
	Invest_type    string    `json:"invest_type"`
	Asset_name     string    `json:"asset_name"`
	Asset_symbol   string    `json:"asset_symbol"`
	Quantity       float64   `json:"quantity"`
	Purchace_price float64   `json:"purchase_price"`
	Current_price  float64   `json:"current_price"`
	Purchase_date  time.Time `json:"purchase_date"`
}

// mutual fund details request type
type mutualfundsreq struct {
	Folio_no string `json:"folio_no"`
}

// response type for mutual funds details
type mutualfundsres struct {
	Demat_no string         `json:"demat_no"`
	Funds    []fundsdetails `json:"funds"`
}

type fundsdetails struct {
	Fund_name      string    `json:"fund_name"`
	Fund_type      string    `json:"fund_type"`
	Quantity       float64   `json:"quantity"`
	Purchace_price float64   `json:"purchase_price"`
	Current_price  float64   `json:"current_price"`
	Purchase_date  time.Time `json:"purchase_date"`
}

// db connection based on database name
func ConnectDB() *sql.DB {
	db_pass := os.Getenv("PG_PASSWORD")
	dsn := fmt.Sprintf("postgres://postgres:%v@localhost/investment?sslmode=disable", db_pass)
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

	router.HandleFunc("/get_investments", get_investments).Methods("POST", "OPTIONS")
	router.HandleFunc("/get_mutualfunds", get_mutualfunds).Methods("POST", "OPTIONS")
	router.HandleFunc("/get_stocks", get_stocks).Methods("POST", "OPTIONS")
	router.HandleFunc("/get_mutualfundsdetail", get_mutualfundsdetails).Methods("POST", "OPTIONS")

	return router
}

// find demat account
func get_investments(w http.ResponseWriter, r *http.Request) {
	var investmentreq investreq

	err := json.NewDecoder(r.Body).Decode(&investmentreq)
	if err != nil {
		http.Error(w, "Enter valid phone number", http.StatusBadRequest)
	}

	db := ConnectDB()
	defer db.Close()
	statement := "select i.name, d.demat_number from demat_accounts as d join investment_brokers as i on i.broker_id=d.broker_id where d.user_id=(select user_id from users where phone_number=$1);"
	row := db.QueryRow(statement, investmentreq.Phone)
	var response investres
	errors := row.Scan(&response.Broker, &response.Demat_no)
	if errors != nil {
		http.Error(w, "account not found", http.StatusNotFound)

	}

	json.NewEncoder(w).Encode(response)

}

// find mutual fund accound(folio)
func get_mutualfunds(w http.ResponseWriter, r *http.Request) {
	var investmentreq investreq

	err := json.NewDecoder(r.Body).Decode(&investmentreq)
	if err != nil {
		http.Error(w, "Enter valid phone number", http.StatusBadRequest)
	}

	db := ConnectDB()
	defer db.Close()
	statement := "select i.name, mf.folio_number from mf_accounts as mf join investment_brokers as i on i.broker_id=mf.broker_id where mf.user_id=(select user_id from users where phone_number=$1);"
	row := db.QueryRow(statement, investmentreq.Phone)
	var response mfresponse
	errors := row.Scan(&response.Broker, &response.Folio_no)
	if errors != nil {
		http.Error(w, "account not found", http.StatusNotFound)

	}

	json.NewEncoder(w).Encode(response)
}

// get investment details
func get_stocks(w http.ResponseWriter, r *http.Request) {
	var stockreq stocksreq

	err := json.NewDecoder(r.Body).Decode(&stockreq)
	if err != nil {
		http.Error(w, "Enter valid demat number", http.StatusBadRequest)
	}

	db := ConnectDB()
	defer db.Close()
	statement := "select investment_type,asset_name,asset_symbol,quantity,purchase_price,current_price,purchase_date from stocks_and_bonds where demat_id=(select demat_id from demat_accounts where demat_number=$1);"
	rows, _ := db.Query(statement, stockreq.Demat_no)
	var response stocksres

	var stocks []stockdetails
	for rows.Next() {
		var stock stockdetails
		errors := rows.Scan(&stock.Invest_type, &stock.Asset_name, &stock.Asset_symbol, &stock.Quantity, &stock.Purchace_price, &stock.Current_price, &stock.Purchase_date)
		if errors != nil {
			http.Error(w, "account not found", http.StatusNotFound)

		}
		stocks = append(stocks, stock)

	}
	response.Demat_no = stockreq.Demat_no
	response.Stocks = stocks
	json.NewEncoder(w).Encode(response)

}

// get mutual fund  details
func get_mutualfundsdetails(w http.ResponseWriter, r *http.Request) {
	var mfreq mutualfundsreq

	err := json.NewDecoder(r.Body).Decode(&mfreq)
	if err != nil {
		http.Error(w, "Enter valid demat number", http.StatusBadRequest)
	}

	db := ConnectDB()
	defer db.Close()
	statement := "select fund_name,fund_type,quantity,purchase_price,current_price,purchase_date from mutual_funds where mf_account_id=(select mf_account_id from mf_accounts where folio_number=$1);"
	rows, _ := db.Query(statement, mfreq.Folio_no)
	var response mutualfundsres

	var funds []fundsdetails
	for rows.Next() {
		var fund fundsdetails
		errors := rows.Scan(&fund.Fund_name, &fund.Fund_type, &fund.Quantity, &fund.Purchace_price, &fund.Current_price, &fund.Purchase_date)
		if errors != nil {
			http.Error(w, "account not found", http.StatusNotFound)

		}
		funds = append(funds, fund)

	}
	response.Demat_no = mfreq.Folio_no
	response.Funds = funds
	json.NewEncoder(w).Encode(response)

}

func main() {
	r := route()
	fmt.Println("Server running at port 8000")
	http.ListenAndServe(":8000", r)
}
