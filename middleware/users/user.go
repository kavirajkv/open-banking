package users

import (
	"encoding/json"
	"fmt"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
	"middleware/database"
	"net/http"
	"os"
	"time"
)

func twilioclient() *twilio.RestClient {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})
	return client
}

func SendOTP(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB("openbanking")
	defer db.Close()
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Enter valid details", http.StatusBadRequest)
		return
	}

	if user.Phone == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}
	var userscheck Users
	usercheck:="SELECT phone_number FROM users WHERE phone_number = $1;"
	er:=db.QueryRow(usercheck,user.Phone).Scan(&userscheck)
	if er!=nil{
		params := &openapi.CreateVerificationParams{}
		params.SetTo(user.Phone)
		params.SetChannel("sms")
		client := twilioclient()
		_, er := client.VerifyV2.CreateVerification(os.Getenv("TWILIO_SERVICE_SID"), params)
		if er != nil {
			http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		}else{
			json.NewEncoder(w).Encode("OTP sent successfully")
		}
		
	}else{
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

}

func VerifyandRegister(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB("openbanking")
	defer db.Close()

	var user AuthUserResgister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Enter valid details", http.StatusBadRequest)
		return
	}

	if user.Phone == "" || user.Otp == "" {
		http.Error(w, "Phone number and OTP are required", http.StatusBadRequest)
		return
	}

	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(user.Phone)
	params.SetCode(user.Otp)

	client := twilioclient()
	_, er := client.VerifyV2.CreateVerificationCheck(os.Getenv("TWILIO_SERVICE_SID"), params)
	if er != nil {
		http.Error(w, "Failed to verify OTP", http.StatusInternalServerError)
		return
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

	json.NewEncoder(w).Encode("User registered successfully")
	
}

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

// func Login(w http.ResponseWriter, r *http.Request) {
// 	db := database.ConnectDB("openbanking")
// 	defer db.Close()

// 	var auth AuthRequest
// 	err := json.NewDecoder(r.Body).Decode(&auth)
// 	if err != nil {
// 		http.Error(w, "Enter valid details", http.StatusBadRequest)
// 	}

// 	statement := "SELECT phone_number FROM users WHERE phone_number = $1;"
// 	row := db.QueryRow(statement, auth.Phone)

// 	var userPhone string
// 	err = row.Scan(&userPhone)
// 	if err != nil {
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(AuthResponse{Phone: userPhone})
// }

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
