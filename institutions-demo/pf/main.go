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


type searchacc struct{
	Phone string `json:"phone"`
}

//nps account find
type npsres struct{
	Pran_no string `json:"pran_no"`
}

//find epf account
type epfres struct{
	Uan_no string `json:"uan_no"`
}

//find ppf account
type ppfres struct{
	Ppf_no string `json:"ppf_no"`
}

//nps account details
type npsdetailsreq struct{
	Pran_no string `json:"pran_no"`
}

type npsdetailsres struct{
	Pran_no string `json:"pran_no"`
	Tier string `json:"tier"`
	Balance float64 `json:"balance"`
	Monthly_contribution time.Time `json:"monthly_contribution"`
	Last_contribution_at time.Time `json:"last_contribution_at"`
}




// db connection based on database name
func ConnectDB() *sql.DB {
	db_pass := os.Getenv("PG_PASSWORD")
	dsn := fmt.Sprintf("postgres://postgres:%v@localhost/providentfund?sslmode=disable", db_pass)
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

	router.HandleFunc("/get_npsaccount", get_npsaccount).Methods("POST", "OPTIONS")
	router.HandleFunc("/get_epfaccount", get_epfaccount).Methods("POST", "OPTIONS")
	router.HandleFunc("/get_ppfaccount", get_ppfaccount).Methods("POST", "OPTIONS")

	return router
}


// get nps account
func get_npsaccount(w http.ResponseWriter, r *http.Request) {
	var npsacc searchacc

	err := json.NewDecoder(r.Body).Decode(&npsacc)
	if err != nil {
		http.Error(w, "Enter valid phone number", http.StatusBadRequest)
	}

	db := ConnectDB()
	defer db.Close()
	statement:="select pran_number from nps_accounts where user_id=(select user_id from users where phone_number=$1);"
	row := db.QueryRow(statement, npsacc.Phone)
	var response npsres

	errors := row.Scan(&response.Pran_no)
	if errors != nil {
		http.Error(w, "account not found", http.StatusNotFound)

	}
	json.NewEncoder(w).Encode(response)

}

// get ppf account
func get_ppfaccount(w http.ResponseWriter, r *http.Request) {
	var ppfacc searchacc

	err := json.NewDecoder(r.Body).Decode(&ppfacc)
	if err != nil {
		http.Error(w, "Enter valid phone number", http.StatusBadRequest)
	}

	db := ConnectDB()
	defer db.Close()
	statement:="select ppf_number from ppf_accounts where user_id=(select user_id from users where phone_number=$1);"
	row := db.QueryRow(statement, ppfacc.Phone)
	var response ppfres

	errors := row.Scan(&response.Ppf_no)
	if errors != nil {
		http.Error(w, "account not found", http.StatusNotFound)

	}
	json.NewEncoder(w).Encode(response)

}

//to get epf account
func get_epfaccount(w http.ResponseWriter, r *http.Request) {
	var epfacc searchacc

	err := json.NewDecoder(r.Body).Decode(&epfacc)
	if err != nil {
		http.Error(w, "Enter valid phone number", http.StatusBadRequest)
	}

	db := ConnectDB()
	defer db.Close()
	statement:="select uan from epf_accounts where user_id=(select user_id from users where phone_number=$1);"
	row := db.QueryRow(statement, epfacc.Phone)
	var response epfres

	errors := row.Scan(&response.Uan_no)
	if errors != nil {
		http.Error(w, "account not found", http.StatusNotFound)

	}
	json.NewEncoder(w).Encode(response)

}


func main() {
	r := route()
	fmt.Println("Server running at port 8000")
	http.ListenAndServe(":8000", r)
}
