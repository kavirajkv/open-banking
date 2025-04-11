package main

import(
	"middleware/routes"
	"fmt"
	"net/http"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {	
	r := routes.Router()
	
	fmt.Println("Server running at port 8080")
	http.ListenAndServe(":8080", r)	
	
}