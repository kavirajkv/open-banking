package users


type Users struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	Aathar string `json:"aathar"`
	Dob    string `json:"dob"`
}

type RegisterResponse struct {
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
	Otp   string `json:"otp"`
}

