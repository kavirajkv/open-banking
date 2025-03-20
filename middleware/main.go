package main

import(
	"middleware/Aggregation"
	"fmt"
	"net/http"
)

func main() {	
	r := Aggregation.Route()
	
	http.ListenAndServe(":8080", r)
	fmt.Println("Server running at port 8080")
	
	
}