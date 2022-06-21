package dto

type ResponseErr struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type BadReqErrResponse struct {
	Message     string      `json:"message"`
	FailedField string      `json:"failed_field"`
	Value       interface{} `json:"value"`
}
