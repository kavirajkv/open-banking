package transfer

// Data request from third party using consent
type DataRequest struct {
	ConsentId       string `json:"consent_id"`
	Phone           string `json:"phone"`
	TpaId           string `json:"tpa_id"`
	InstitutionId   string `json:"institution_id"`
	DataType        string `json:"data_type"`
	DataLimit       string `json:"data_limit"`
}

// data response with public key , data and nonce and signature
type DataResponse struct {
	ConsentId     string `json:"consent_id"`
	Phone         string `json:"phone"`
	InstitutionId string `json:"institution_id"`
	ResponseData  string `json:"response_data"`
	Status        string `json:"status"`
	PublicKey     string `json:"public_key"`
	Nonce         string `json:"nonce"`
}

