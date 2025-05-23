package consent

import (
	"encoding/json"
	"fmt"
	"middleware/database"
	"net/http"
	"time"
	"github.com/lib/pq"
	"github.com/google/uuid"
)

// func to create consent
func CreateConsent(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB("openbanking")
	defer db.Close()

	var consent Consent
	err := json.NewDecoder(r.Body).Decode(&consent)
	if err != nil {
		http.Error(w, "Enter valid details", http.StatusBadRequest)
	}

	consentid := "consent_"+uuid.New().String()
	expiry, err := time.Parse("2006-01-02", consent.Expiry)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}
	status:="Active"

	statement := "INSERT INTO consent (consent_id,user_did,tpa_id,financial_institution_id,requested_data,expiry_time,status) VALUES ($1,$2,$3,$4,$5,$6,$7);"
	res, err := db.Exec(statement, consentid, consent.User_did, consent.Tpa_id, consent.Financial_institution_id, pq.Array(consent.Requested_data), expiry, status)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected == 0 {
		http.Error(w, "Failed to create consent", http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Insertion successful for %v", rowsAffected)

	response := ConsentResponse{Message: msg}

	json.NewEncoder(w).Encode(response)
}


// func to get consent details
func GetConsentDetails(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB("openbanking")
	defer db.Close()

	var req ConsentDetailsrequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Enter valid details", http.StatusBadRequest)
	}

	statement := "SELECT consent_id,tpa_id,financial_institution_id,requested_data,expiry_time,status FROM consent WHERE user_did=(select did from did where user_id=(select user_id from users where phone_number=$1));"
	rows,err:= db.Query(statement, req.Phone)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var consents []ConsentDetailsResponse

	for rows.Next() {
		var consent ConsentDetailsResponse
		err = rows.Scan(&consent.Consent_id,&consent.Tpa_id,&consent.Financial_institution_id,pq.Array(&consent.Requested_data),&consent.Expiry,&consent.Status)
		if err != nil {
			http.Error(w, "Consent not found", http.StatusNotFound)
			return
		}
		consents = append(consents, consent)
	}

	json.NewEncoder(w).Encode(consents)
}


func GetConsentbyId(w http.ResponseWriter,r *http.Request){
	db := database.ConnectDB("openbanking")
	defer db.Close()

	var req ConsentbyIdreq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Enter valid details", http.StatusBadRequest)
	}

	statement := "SELECT user_did,tpa_id,financial_institution_id,requested_data,expiry_time,status FROM consent WHERE consent_id=$1;"
	row:= db.QueryRow(statement, req.Consent_id)
	var consent ConsentbyIdres
	err = row.Scan(&consent.User_did,&consent.Tpa_id,&consent.Financial_institution_id,pq.Array(&consent.Requested_data),&consent.Expiry,&consent.Status)
	if err != nil {
		http.Error(w, "Consent not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(consent)
}


// func to update consent
func UpdateConsent(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB("openbanking")
	defer db.Close()

	var req UpdateConsentStatusreq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Enter valid details", http.StatusBadRequest)
	}

	statement := "UPDATE consent SET status=$1 WHERE consent_id=$2;"
	res, err := db.Exec(statement, req.Status, req.Consent_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected == 0 {
		http.Error(w, "Failed to update consent", http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Update successful for %v", rowsAffected)

	response := ConsentResponse{Message: msg}

	json.NewEncoder(w).Encode(response)
}