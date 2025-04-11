package Aggregation

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
)


// get all bank account details
func AggregateBankAccount(w http.ResponseWriter, r *http.Request) {
	var accountReq Account_req

	err := json.NewDecoder(r.Body).Decode(&accountReq)
	if err != nil {
		http.Error(w, "Enter valid phone number", http.StatusBadRequest)
	}

	// financial institution API endpoints
	endpoints := map[string]string{
		"icici": "http://127.0.0.1:8000/icici/get_account",
		"axis":  "http://127.0.0.1:8000/axis/get_account",
		"hdfc":  "http://127.0.0.1:8000/hdfc/get_account",
	}

	requestBody, err := json.Marshal(accountReq)
	if err != nil {
		http.Error(w, "Failed to marshal phone number", http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup
	responses := make([]InstitutionResponse, 0, len(endpoints))
	for institute, endpoint := range endpoints {
		wg.Add(1)
		go func(institute string, endpoint string) {
			defer wg.Done()

			resp, err := http.Post(endpoint, "application/json", bytes.NewReader(requestBody))
			if err != nil {
				responses = append(responses, InstitutionResponse{
					Institution: institute,
					Error:       err.Error(),
				})
				return
			}
			defer resp.Body.Close()

			// Read and parse response
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				responses = append(responses, InstitutionResponse{
					Institution: institute,
					Error:       err.Error(),
				})
				return
			}

			//interface to hold response from financial institution
			var data interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				responses = append(responses, InstitutionResponse{
					Institution: institute,
					Error:       err.Error(),
				})
				return
			}

			responses = append(responses, InstitutionResponse{
				Institution: institute,
				Data:        data,
			})
		}(institute, endpoint)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	json.NewEncoder(w).Encode(responses)
}


//func to get data from specified institute of an user
func GetData(w http.ResponseWriter,r *http.Request){
	
}
