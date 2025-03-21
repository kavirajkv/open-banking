package users

import (
	"encoding/json"
	"fmt"
	// "log"
	"middleware/database"
	"net/http"
	"time"
)

// func to register user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB("openbanking")
	defer db.Close()

	var user Users
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Enter valid details", http.StatusBadRequest)
	}

	// Convert DOB string to time.Time
	dob, err := time.Parse("1999-01-02", user.Dob)
	if err != nil {
		http.Error(w, "Invalid date format for DOB", http.StatusBadRequest)
		return
	}

	statement := "INSERT INTO users (phone_number, name, aadhar, dob, email) VALUES ($1, $2, $3, $4, $5);"
	res, err := db.Exec(statement, user.Phone, user.Name, user.Aathar, dob, user.Email)

	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected == 0 {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Insertion successful for %v", rowsAffected)

	response := RegisterResponse{Message: msg}

	json.NewEncoder(w).Encode(response)
}

// func to get user details
func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB("openbanking")
	defer db.Close()

	var phone DetailsRequest
	err := json.NewDecoder(r.Body).Decode(&phone)
	if err != nil {
		http.Error(w, "Enter valid phone number", http.StatusBadRequest)
	}

	statement := "SELECT name, phone_number, aadhar, dob, email FROM users WHERE phone_number = $1;"
	row := db.QueryRow(statement, phone.Phone)

	var userDetails Users
	err = row.Scan(&userDetails.Name, &userDetails.Phone, &userDetails.Aathar, &userDetails.Dob, &userDetails.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(userDetails)
}
