package dto

type CheckinVerifyRequest struct {
	EventType int `json:"event_type" example:"1"`
}

type CheckinVerifyResponse struct {
	CheckinToken string `json:"checkin_token" example:"ec5b9355-0b6c-11ed-b88b-0250cf8509e4"`
	CheckinType  int32  `json:"checkin_type" example:"1"`
}

type CheckinConfirmRequest struct {
	Token string `json:"token"`
}
