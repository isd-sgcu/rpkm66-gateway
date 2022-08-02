package dto

type VerifyEstampRequest struct {
	EventId string `json:"event_id" example:"ec5b9355-0b6c-11ed-b88b-0250cf8509e4"`
}

type ConfirmEstampRequest struct {
	EventId string `json:"event_id" example:"ec5b9355-0b6c-11ed-b88b-0250cf8509e4"`
}
