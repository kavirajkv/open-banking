package users

import (
	"github.com/golang-jwt/jwt/v5"
)

type Users struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	Aathar string `json:"aathar"`
	Dob    string `json:"dob"`
}

type Response struct {
	Message string `json:"message"`
}

type DetailsRequest struct {
	Phone string `json:"phone"`
}

type User struct {
	Phone string `json:"phone"`
}

type AuthUserResgister struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	Aathar string `json:"aathar"`
	Dob    string `json:"dob"`
	Otp    string `json:"otp"`
}

type LoginRequest struct {
	Phone string `json:"phone"`
	Otp   string `json:"otp"`
}

type Claims struct {
	Phone string `json:"phone"`
	jwt.RegisteredClaims
}
