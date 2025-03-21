package routes

import (
	"middleware/Aggregation"
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

	

	return router
}