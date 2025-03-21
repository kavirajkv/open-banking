package main

import(
	"middleware/routes"
	"fmt"
	"net/http"
)

func main() {	
	r := routes.Router()
	
	fmt.Println("Server running at port 8080")
	http.ListenAndServe(":8080", r)	
	
}