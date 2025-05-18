package nominee



type Nomineereq struct{
	Acc_no string `json:"acc_no"`
	Bank string `json:"bank"`
}

type Nomineeres struct{
	NomineeName string `json:"nominee"`
	Relation string `json:"relation"`
	Status string `json:"status"`
}

