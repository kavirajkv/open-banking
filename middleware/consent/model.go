package consent

type Consent struct {
	User_did                 string   `json:"user_did"`
	Tpa_id                   string   `json:"tpa_id"`
	Financial_institution_id string   `json:"provider_id"`
	Requested_data           []string `json:"requested_data"`
	Expiry                   string   `json:"expiry"`
}

type ConsentResponse struct {
	Message string `json:"message"`
}
