package routes

import (
	"middleware/Aggregation"
	"middleware/consent"
	"middleware/users"

	"github.com/gorilla/mux"
)

func Router() *mux.Router{
	router:=mux.NewRouter()

	//  Aggregation routes  //
	//route to get all bank account details
	router.HandleFunc("/get_aggregated_accounts", Aggregation.AggregateBankAccount).Methods("POST", "OPTIONS")


	// User routes  //
	//user registration route
	router.HandleFunc("/register_user", users.RegisterUser).Methods("POST", "OPTIONS")
	//get user details route
	router.HandleFunc("/get_user_details", users.GetUserDetails).Methods("POST", "OPTIONS")

	// Consent routes //
	//create consent
	router.HandleFunc("/create_consent",consent.CreateConsent).Methods("POST","OPTIONS")
	//get consent details
	router.HandleFunc("/get_consent_details",consent.GetConsentDetails).Methods("POST","OPTIONS")

	return router
}