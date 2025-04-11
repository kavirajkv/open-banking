package Aggregation

// Account_req struct to hold phone number
type Account_req struct {
	Phone string `json:"phone"`
}

// InstitutionResponse struct to hold response from financial institution
type InstitutionResponse struct {
	Institution string      `json:"institution"`
	Data        interface{} `json:"data"`
	Error       string      `json:"error"`
}

type DataRequest struct {
	Phone string `json:"phone"`
	Consent_id string `json:"consent_id"`
}

type DataResponse struct {
	Data interface{} `json:"data"`
	Error string `json:"error"`
}

