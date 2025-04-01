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

type ConsentDetailsrequest struct {
	User_did string `json:"user_did"`
}

type ConsentDetailsResponse struct {
	Consent_id               string   `json:"consent_id"`
	Tpa_id                   string   `json:"tpa_id"`
	Financial_institution_id string   `json:"provider_id"`
	Requested_data           []string `json:"requested_data"`
	Expiry                   string   `json:"expiry"`
	Status                   string   `json:"status"`
}

type ConsentbyIdreq struct {
	Consent_id string `json:"consent_id"`
}

type ConsentbyIdres struct {
	User_did string `json:"user_did"`
	Tpa_id                   string   `json:"tpa_id"`
	Financial_institution_id string   `json:"provider_id"`
	Requested_data           []string `json:"requested_data"`
	Expiry                   string   `json:"expiry"`
	Status                   string   `json:"status"`
}


type UpdateConsentStatusreq struct {
	Consent_id string `json:"consent_id"`
	Status     string `json:"status"`
}

