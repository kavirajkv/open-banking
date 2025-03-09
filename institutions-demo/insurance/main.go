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

// insurance find request
type insurancereq struct {
	Phone string `json:"phone"`
}

// insurance find response
type insuranceres struct {
	Policies []insurance `json:"policies"`
}

type insurance struct {
	Provider  string `json:"provider"`
	Policy_no string `json:"policy_no"`
}

// insurance details request type
type insuredetailsreq struct {
	Policy_no string `json:"policy_no"`
}

// insurance details response
type insuredetailsres struct {
	Policy_no string        `json:"policy_no"`
	Details   policydetails `json:"details"`
}

type policydetails struct {
	Provider         string    `json:"provider"`
	Policy_type      string    `json:"policy_type"`
	Sum_assured      float64   `json:"sum_assured"`
	Premium          float64   `json:"premium"`
	Start_date       time.Time `json:"start_date"`
	End_date         time.Time `json:"end_date"`
	Nominee          string    `json:"nominee"`
	Nominee_relation string    `json:"nominee_relation"`
}

// db connection based on database name
func ConnectDB() *sql.DB {
	db_pass := os.Getenv("PG_PASSWORD")
	dsn := fmt.Sprintf("postgres://postgres:%v@localhost/insurance?sslmode=disable", db_pass)
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

	router.HandleFunc("/get_insurance", get_insurance).Methods("POST", "OPTIONS")
	router.HandleFunc("/get_insurancedetails", get_insurancedetails).Methods("POST", "OPTIONS")

	return router
}

// find insurance policies
func get_insurance(w http.ResponseWriter, r *http.Request) {
	var insurereq insurancereq

	err := json.NewDecoder(r.Body).Decode(&insurereq)
	if err != nil {
		http.Error(w, "Enter valid phone number", http.StatusBadRequest)
	}

	db := ConnectDB()
	defer db.Close()
	statement := "select ip.name,p.policy_number from insurance_policies as p join insurance_providers as ip on ip.provider_id=p.provider_id where user_id=(select user_id from users where phone_number=$1);"
	rows, _ := db.Query(statement, insurereq.Phone)
	var response insuranceres

	var policies []insurance
	for rows.Next() {
		var policy insurance
		errors := rows.Scan(&policy.Provider, &policy.Policy_no)
		if errors != nil {
			http.Error(w, "account not found", http.StatusNotFound)

		}
		policies = append(policies, policy)

	}

	response.Policies = policies
	json.NewEncoder(w).Encode(response)

}

// get insurance policy details
func get_insurancedetails(w http.ResponseWriter, r *http.Request) {
	var insuredetail insuredetailsreq

	err := json.NewDecoder(r.Body).Decode(&insuredetail)
	if err != nil {
		http.Error(w, "Enter valid phone number", http.StatusBadRequest)
	}

	db := ConnectDB()
	defer db.Close()
	statement := "select ip.name,p.policy_type,p.sum_assured,p.premium_amount,p.start_date,p.end_date,p.nominee,p.nominee_relationship from insurance_policies as p join insurance_providers as ip on ip.provider_id=p.provider_id where policy_number=$1;"
	row := db.QueryRow(statement, insuredetail.Policy_no)
	var response insuredetailsres

	errors := row.Scan(&response.Details.Provider, &response.Details.Policy_type, &response.Details.Sum_assured, &response.Details.Premium, &response.Details.Start_date, &response.Details.End_date, &response.Details.Nominee, &response.Details.Nominee_relation)
	if errors != nil {
		http.Error(w, "account not found", http.StatusNotFound)

	}

	response.Policy_no = insuredetail.Policy_no
	json.NewEncoder(w).Encode(response)

}

func main() {
	r := route()
	fmt.Println("Server running at port 8000")
	http.ListenAndServe(":8000", r)
}
