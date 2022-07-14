package dto

type VaccineResponse struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	IsPassed  bool   `json:"is_passed"`
	Uid       string `json:"uid"`
}

type VaccineRequest struct {
	HCert     string `json:"hcert"`
	StudentId string `json:"uid"`
}

type Verify struct {
	HCert string `json:"hcert" validate:"required"`
}
