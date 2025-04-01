package consent

import (
	"encoding/json"
	"fmt"
	"middleware/database"
	"net/http"
	"time"
	"github.com/lib/pq"
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

	consentid := "consent" + time.Now().String() + consent.User_did
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
