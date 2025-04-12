package users

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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

func sendotp(phone string, channel string) error {
	params := &openapi.CreateVerificationParams{}
	params.SetTo(phone)
	params.SetChannel(channel)
	client := twilioclient()
	_, err := client.VerifyV2.CreateVerification(os.Getenv("TWILIO_SERVICE_SID"), params)
	if err != nil {
		return err
	}
	return nil
}

func verifyotp(phone string, otp string) (bool, error) {
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(phone)
	params.SetCode(otp)

	client := twilioclient()

	valid, err := client.VerifyV2.CreateVerificationCheck(os.Getenv("TWILIO_SERVICE_SID"), params)
	if err != nil {
		return false, err
	}
	return *valid.Status == "approved", nil
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
	usercheck := "SELECT phone_number FROM users WHERE phone_number = $1;"
	er := db.QueryRow(usercheck, user.Phone).Scan(&userscheck)
	if er != nil {
		// Send OTP
		er := sendotp(user.Phone, "sms")
		if er != nil {
			http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		} else {
			json.NewEncoder(w).Encode("OTP sent successfully")
		}

	} else {
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

	// Verify OTP
	verify, er := verifyotp(user.Phone, user.Otp)

	if er != nil && verify != true {
		http.Error(w, er.Error(), http.StatusInternalServerError)
		return
	}

	// Convert DOB string to time.Time
	dob, err := time.Parse("2006-01-02", user.Dob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statement := "INSERT INTO users (phone_number, name, aadhar, dob, email) VALUES ($1, $2, $3, $4, $5);"
	res, err := db.Exec(statement, user.Phone, user.Name, user.Aathar, dob, user.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Insertion successful for %v", rowsAffected)

	response := Response{Message: msg}

	json.NewEncoder(w).Encode(response)

}

// to login user
func LoginUserOTP(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB("openbanking")
	defer db.Close()

	var loginuser User

	err := json.NewDecoder(r.Body).Decode(&loginuser)
	if err != nil {
		http.Error(w, "Enter valid phone number", http.StatusBadRequest)
		return
	}

	var userscheck User
	usercheck := "SELECT phone_number FROM users WHERE phone_number = $1;"
	er := db.QueryRow(usercheck, loginuser.Phone).Scan(&userscheck.Phone)

	if er != nil {
		http.Error(w, er.Error(), http.StatusBadRequest)
	}

	//send otp
	otperr := sendotp(loginuser.Phone, "sms")
	if otperr != nil {
		http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("OTP sent successfully")

}

// func to verify user
func LoginOTPverify(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB("openbanking")
	defer db.Close()

	var loginuser LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginuser)
	if err != nil {
		http.Error(w, "Enter valid phone number", http.StatusBadRequest)
	}

	if loginuser.Phone == "" || loginuser.Otp == "" {
		http.Error(w, "Phone number and OTP are required", http.StatusBadRequest)
		return
	}

	//verify otp
	verify, er := verifyotp(loginuser.Phone, loginuser.Otp)

	if er != nil && verify != true {
		http.Error(w, "Failed to verify OTP", http.StatusInternalServerError)
		return
	}

	//here after valid credentials, generate the jwt token and send it to the user (token is set to valid for 10 minutes)
	claim := Claims{Phone: loginuser.Phone,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10))},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(time.Minute * 10),
	})

	msg := Response{Message: "User logged in successfully"}

	json.NewEncoder(w).Encode(msg)
}

// to authenticate user based on token
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}

		tokenstr := cookie.Value

		claim := &Claims{}

		token, err := jwt.ParseWithClaims(tokenstr, claim, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		//here after authencation username is send to the next route
		// ctx:=context.WithValue(r.Context(),"username",claim.Username)
		next(w, r)

	}
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
