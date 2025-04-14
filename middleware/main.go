package main

import(
	"middleware/routes"
	"fmt"
	"net/http"
	"github.com/joho/godotenv"
	"log"
	"github.com/rs/cors"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {	
	r := routes.Router()

	c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"}, // Allows all origins
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
    })

	handler := c.Handler(r)
	
	fmt.Println("Server running at port 8080")
	http.ListenAndServe(":8080", handler)	
	
}